#!/usr/bin/env python3

import argparse
import datetime

from cryptography import x509
from cryptography.hazmat.primitives import hashes, serialization
from cryptography.hazmat.primitives.asymmetric import ec, rsa
from cryptography.hazmat.primitives.serialization import pkcs12
from cryptography.x509.oid import NameOID


def save_key(key: ec.EllipticCurvePrivateKey, file_path: str, password: bytes):
    # Get unencrypted PEM private key description
    encryption = (
        serialization.BestAvailableEncryption(password)
        if password is not None
        else serialization.NoEncryption()
    )

    key_bytes = key.private_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PrivateFormat.PKCS8,
        encryption_algorithm=encryption,
    )

    with open(file_path, "wb") as f:
        f.write(key_bytes)


def load_key(file_path: str, password: bytes) -> ec.EllipticCurvePrivateKey:
    with open(file_path, "rb") as f:
        key_bytes = f.read()

    key = serialization.load_pem_private_key(key_bytes, password)
    return key


def save_cert(
    cert: x509.Certificate, file_path: str, root_cert: x509.Certificate = None
):
    cert_bytes = cert.public_bytes(serialization.Encoding.PEM)

    with open(file_path, "wb") as f:
        f.write(cert_bytes)

    if root_cert is not None:
        root_cert_bytes = root_cert.public_bytes(serialization.Encoding.PEM)
        with open(file_path, "ab") as f:
            f.write(root_cert_bytes)


def load_cert(file_path: str) -> x509.Certificate:
    with open(file_path, "rb") as f:
        cert_bytes = f.read()

    cert = x509.load_pem_x509_certificate(cert_bytes)
    return cert


def generate_private_key() -> ec.EllipticCurvePrivateKey:
    # Use ECC curve NIST P-256
    key = ec.generate_private_key(curve=ec.SECP256R1)
    # key = rsa.generate_private_key(
    #     public_exponent=65537,
    #     key_size=4096,
    # )

    return key


def save_cert_key(
    cert: x509.Certificate,
    key: ec.EllipticCurvePrivateKey,
    name: str,
    cert_store: str,
    password: bytes,
    root_cert: x509.Certificate = None,
):
    # encryption = serialization.BestAvailableEncryption(password)
    # name_bytes = name.encode("utf-8")
    # data = pkcs12.serialize_key_and_certificates(
    #     name_bytes, key, cert, None, encryption
    # )

    # with open("{}/{}.p12".format(cert_store, name), "wb") as f:
    #     f.write(data)
    save_key(key, "{}/{}.key".format(cert_store, name), password)
    save_cert(cert, "{}/{}.crt".format(cert_store, name), root_cert)


def load_cert_key(name: str, cert_store: str, password: bytes):
    # with open("{}/{}.p12".format(cert_store, name), "rb") as f:
    #     data = f.read()

    # key, cert, _ = pkcs12.load_key_and_certificates(data, password)

    key = load_key("{}/{}.key".format(cert_store, name), password)
    cert = load_cert("{}/{}.crt".format(cert_store, name))

    return cert, key


class Entity:
    def __init__(self, country, region, locality, organization, name=None):
        self.country = country
        self.region = region
        self.locality = locality
        self.organization = organization
        self.name = name
        self.url = None

    def add_url(self, url):
        self.url = url

    def x509_name(self) -> x509.Name:
        attributes = [
            x509.NameAttribute(NameOID.COUNTRY_NAME, self.country),
            x509.NameAttribute(NameOID.STATE_OR_PROVINCE_NAME, self.region),
            x509.NameAttribute(NameOID.LOCALITY_NAME, self.locality),
            x509.NameAttribute(NameOID.ORGANIZATION_NAME, self.organization),
        ]

        common_name = self.url or self.name
        if common_name is not None:
            attributes.append(x509.NameAttribute(NameOID.COMMON_NAME, common_name))

        return x509.Name(attributes)

    def x509_sans(self) -> x509.SubjectAlternativeName:
        names = [x509.DNSName(self.url)] if self.url is not None else []

        return x509.SubjectAlternativeName(names)


def generate_cert(
    private_key: ec.EllipticCurvePrivateKey,
    signing_private_key: ec.EllipticCurvePrivateKey,
    issuer_skid: x509.SubjectKeyIdentifier,
    subject: Entity,
    issuer: x509.Name,
    expiration_days: int,
    is_ca=False,
):
    public_key = private_key.public_key()

    builder = x509.CertificateBuilder()
    builder = builder.public_key(public_key)
    builder = builder.serial_number(x509.random_serial_number())

    builder = builder.subject_name(subject.x509_name())
    builder = builder.issuer_name(issuer)

    # Assign expiration date
    expiration_date = datetime.datetime.utcnow() + datetime.timedelta(
        expiration_days, 0, 0
    )
    builder = builder.not_valid_before(datetime.datetime.utcnow())
    builder = builder.not_valid_after(expiration_date)

    # Critical extensions cause certificate to be rejected if not supported
    # builder = builder.add_extension(subject.x509_sans(), critical=False)
    # Ensures this CA cannot create subordinate CAs - unnecessary for just generating robot/server certs
    path_length = 0 if is_ca else None
    ca_constraints = x509.BasicConstraints(ca=is_ca, path_length=path_length)
    builder = builder.add_extension(ca_constraints, critical=True)
    skid = x509.SubjectKeyIdentifier.from_public_key(public_key)
    builder = builder.add_extension(skid, critical=False)
    if not is_ca:
        akid = x509.AuthorityKeyIdentifier.from_issuer_subject_key_identifier(
            issuer_skid
        )
        builder = builder.add_extension(akid, critical=False)

    # SHA256 is better supported and leads to smaller digests than SHA512
    certificate = builder.sign(
        private_key=signing_private_key, algorithm=hashes.SHA256()
    )

    return certificate


def create_root(args, ca_entity):
    private_key = generate_private_key()

    cert = generate_cert(
        private_key,
        private_key,
        None,
        ca_entity,
        ca_entity.x509_name(),
        args.expire,
        is_ca=True,
    )

    return cert, private_key, None


def issue_cert(args, entity):
    if args.url:
        entity.add_url(args.url)

    ca_password = args.ca_password.encode("utf-8")
    ca_password = None if len(ca_password) == 0 else ca_password
    ca_cert, ca_private_key = load_cert_key(args.ca, args.store, ca_password)
    private_key = generate_private_key()
    ca_skid = ca_cert.extensions.get_extension_for_class(x509.SubjectKeyIdentifier)

    cert = generate_cert(
        private_key,
        ca_private_key,
        ca_skid.value,
        entity,
        ca_cert.subject,
        args.expire,
        is_ca=False,
    )

    return cert, private_key, ca_cert


def main():
    parser = argparse.ArgumentParser(
        formatter_class=argparse.ArgumentDefaultsHelpFormatter
    )
    subparsers = parser.add_subparsers(title="subcommands", help="see additional help")
    ca_parser = subparsers.add_parser(
        "create_ca",
        help="create root CA",
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )
    ca_parser.set_defaults(func=create_root)
    issue_parser = subparsers.add_parser(
        "issue",
        help="issue cert from CA",
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )
    issue_parser.set_defaults(func=issue_cert)

    parser.add_argument(
        "-l",
        "--location",
        type=str,
        nargs=3,
        default=["US", "Maryland", "Baltimore"],
        help="subject location",
    )
    parser.add_argument(
        "-o", "--org", type=str, default="Deliverbot", help="organization name"
    )
    parser.add_argument(
        "-s", "--store", type=str, default="./", help="where to store certs and keys"
    )
    parser.add_argument(
        "-p", "--password", type=str, default="", help="key encryption password"
    )

    ca_parser.add_argument(
        "-n", "--name", type=str, required=True, help="CA name, e.g. deliverbot_ca"
    )
    ca_parser.add_argument(
        "-e", "--expire", type=int, default=30, help="days until CA expires"
    )
    issue_parser.add_argument(
        "-n",
        "--name",
        type=str,
        required=True,
        help="subject name, e.g. deliverbot_server",
    )
    issue_parser.add_argument(
        "-u", "--url", type=str, required=False, help="URL for server certs"
    )
    issue_parser.add_argument(
        "-c", "--ca", type=str, default="", help="CA name to issue cert from"
    )
    issue_parser.add_argument(
        "--ca_password", type=str, default="", help="password for CA"
    )
    issue_parser.add_argument(
        "-e", "--expire", type=int, default=30, help="days until cert expires"
    )

    args = parser.parse_args()

    entity = Entity(
        args.location[0],
        args.location[1],
        args.location[2],
        args.org,
        args.name,
    )
    cert, key, root = args.func(args, entity)
    password = args.password.encode("utf-8") if len(args.password) > 0 else None
    save_cert_key(cert, key, args.name, args.store, password, root)


if __name__ == "__main__":
    main()
