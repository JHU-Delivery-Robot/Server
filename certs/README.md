# Certificates

We rely on mTLS both to encryt communications between the server and the robots, as well as for authentication. All security is at the transport level: if you aren't authorized to use the server, you simply cannot connect at all. There are no passwords or other tokens - having an appropriate certificate provides proof of identity.

certgen.py is a tool for generating mTLS certificates

Setting up mTLS security involves two steps
 - generating a root CA (certificate authority) which will issue all the mTLS certificates.
 - using the CA to issue certificates to the server and robots

Each certificate (including the root CA) has an accompaning private key, which must be kept secret and be password protected.

To generating a root CA named "deliverbot_ca":
```
python .\certgen.py -p <ca password> ca -n deliverbot_ca
```

To issue a cert for a hypothetical server named "deliverbot_server" running at `fangornsbane.com`:
```
python .\certgen.py -p <cert password> issue -n deliverbot_server -c deliverbot_ca --ca_password <ca password> -u fangornsbane.com
```

The certgen utility also allows changing the expiration period (default is 30 days) and modifying the subject organization info for the certificates.

The CA certificate needs to be set as the root CA on both the server and the robots. The CA key file and password should be kept secret, they are needed to issue additional certificates.
