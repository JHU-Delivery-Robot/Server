python .\certgen.py -p fangorn ca -n deliverbot_ca
python .\certgen.py -p bane issue -n deliverbot_server -c deliverbot_ca --ca_password fangorn -u localhost
python .\certgen.py -p bane issue -n deliverbot_robot -c deliverbot_ca --ca_password fangorn
