## README


`fakerize.py` is a very simple script is designed to create synthetic records for testing purposes, since true records contain sensitive information such as cryptographic signatures and PII. Currently the script is able to create synthetic VASP records and certificate requests, but could be expanded to generate other kinds of data.

The records are synthesized using templates available in the `templates` folder as well as some 3rd party Python libraries [`Faker`](https://faker.readthedocs.io/en/master/) and [`lorem`](https://pypi.org/project/lorem/). Ensure you have install the requirements with `pip install -r requirements.txt`. Then run `python fakerize.py`. To modify the features of the synthetic VASPS (or certificate requests), such as adding a new state, modify the `FAKE_VASPS` (or `FAKE_CERTS`) and add the new state to `VASP_STATE_CHANGES` (or `CERT_STATE_CHANGES`).
