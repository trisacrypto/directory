# fakerize.py
# Clean and fakerize the contents of the GDS database to produce safe test data
#
#
# Author:   Rebecca Bilbro
# Created:  Sun Oct 31 13:20:21 EDT 2021
#
#
##########################################################################
# Imports
##########################################################################

import os
import json
import uuid
import lorem
import random
import secrets
import datetime
from faker import Faker


##########################################################################
# Global Variables
##########################################################################

COUNTRIES = ["US", "CA", "CN", "CI", "GR", "GY", "MG", "MA", "DE", "SG"]
DOMAINS = [".com", ".net", ".io", ".ai", ".org"]
VASP_CATEGORIES = [
    "Exchange",
    "DEX",
    "P2P",
    "Kiosk",
    "Custodian",
    "OTC",
    "Fund",
    "Project",
    "Gambling",
    "Miner",
    "Mixer",
    "Individual",
    "Other",
]
FAKE_VASPS = {
    "CharlieBank": "SUBMITTED",
    "Delta Assets": "APPEALED",
    "Echo Funds": "SUBMITTED",
    "Foxtrot LLC": "VERIFIED",
    "GolfBucks": "ERRORED",
    "Hotel Corp": "VERIFIED",
    "IndiaCoin": "VERIFIED",
    "Juliet Capulet LLC": "PENDING_REVIEW",
    "KiloVASP": "VERIFIED",
    "Lima Beancounters": "REJECTED",
    "Mikes Official VASP": "REJECTED",
    "NovemberCash": "VERIFIED",
    "Romeo Montague Labs LLC": "VERIFIED",
    "Oscar Inc": "PENDING_REVIEW",
}
STATE_CHANGES = {
    "SUBMITTED": {
        "previous_state": "NO_VERIFICATION",
        "current_state": "SUBMITTED",
        "description": "register request received",
        "source": "automated",
    },
    "ERRORED": {
        "previous_state": "SUBMITTED",
        "current_state": "ERRORED",
        "description": "registration request error",
        "source": "automated",
    },
    "EMAIL_VERIFIED": {
        "previous_state": "SUBMITTED",
        "current_state": "EMAIL_VERIFIED",
        "description": "completed email verification",
        "source": "automated",
    },
    "PENDING_REVIEW": {
        "previous_state": "EMAIL_VERIFIED",
        "current_state": "PENDING_REVIEW",
        "description": "review email sent",
        "source": "automated",
    },
    "REVIEWED": {
        "previous_state": "PENDING_REVIEW",
        "current_state": "REVIEWED",
        "description": "registration request received",
        "source": "admin@rotational.io",
    },
    "ISSUING_CERTIFICATE": {
        "previous_state": "REVIEWED",
        "current_state": "ISSUING_CERTIFICATE",
        "description": "issuing certificate",
        "source": "automated",
    },
    "VERIFIED": {
        "previous_state": "ISSUING_CERTIFICATE",
        "current_state": "VERIFIED",
        "description": "certificate issued",
        "source": "automated",
    },
    "REJECTED": {
        "previous_state": "PENDING_REVIEW",
        "current_state": "REJECTED",
        "description": "registration rejected",
        "source": "admin@rotational.io",
    },
    "APPEALED": {
        "previous_state": "REJECTED",
        "current_state": "APPEALED",
        "description": "registration appealed",
        "source": "admin@rotational.io",
    },
}

##########################################################################
# Helper Methods - Used by subsequent functions to synthesize VASP records
##########################################################################


def fake_legal_name(vasp):
    """
    Given a string representing a VASP's name, return a valid
    dictionary for the faked name identifiers
    """
    return {
        "name_identifiers": [
            {
                "legal_person_name": vasp,
                "legal_person_name_identifier_type": "LEGAL_PERSON_NAME_TYPE_CODE_LEGL",
            }
        ]
    }


def fake_address(country):
    """
    Given a string representing the 2-letter country code for a VASP, return a
    dictionary for the faked geographic address
    """
    fake = Faker()
    return {
        "address_type": "ADDRESS_TYPE_CODE_BIZZ",
        "department": "",
        "sub_department": "",
        "street_name": "",
        "building_number": "",
        "building_name": "",
        "floor": "",
        "post_box": "",
        "room": "",
        "post_code": "",
        "town_name": "",
        "town_location_name": "",
        "district_name": "",
        "country_sub_division": "",
        "address_line": [fake.street_address()],
        "country": country,
    }


def make_person(vasp):
    """
    Given a string representing a VASP's name, return a valid
    dictionary for a representative of that VASP including faked name
    and contact information.
    """
    fake = Faker()
    name = fake.name()
    domain = random.choice(DOMAINS)
    email = name.lower().split()[0] + "@" + vasp.replace(" ", "").lower() + domain
    return {
        "name": name,
        "email": email,
        "phone": fake.phone_number(),
        "extra": {
            "@type": "type.googleapis.com/gds.models.v1.GDSContactExtraData",
            "verified": True,
            "token": "",
            "email_log": [],
        },
    }


def make_dates(first="2021-06-15T05:11:13Z", last="2021-10-25T17:15:43Z", count=3):
    """
    Make `count` number of sequential dates between `first` and `last`,
    return the dates as a list.
    """
    format = "%Y-%m-%dT%H:%M:%SZ"
    start = datetime.datetime.strptime(first, format)
    end = datetime.datetime.strptime(last, format)
    dates = [random.random() * (end - start) + start for _ in range(count)]
    return [datetime.datetime.strftime(date, format) for date in sorted(dates)]


def make_log(state="VERIFIED"):
    """
    Create a fake audit log depending on the dates provided
    and the setting of `state`.
    Returns a list of dictionaries.
    """
    logs = []

    states = [state]
    current_state = state
    prior_state = None
    while current_state in STATE_CHANGES:
        prior_state = STATE_CHANGES[current_state]["previous_state"]
        states.insert(0, prior_state)
        current_state = prior_state

    # skip the first data and state, since it's NO_VERIFICATION and doesn't get a log
    dates = make_dates(count=(len(states) - 1))

    for st, dt in zip(states[1:], dates):
        log = dict()
        log.update(STATE_CHANGES[st])
        log["timestamp"] = dt
        logs.append(log)

    return logs


def make_notes():
    """
    Make fake review notes. Returns a nested dictionary, since notes are a dict
    not a list.
    """
    idx = str(uuid.uuid1())
    created, modified = make_dates(count=2)
    return {
        idx: {
            "id": idx,
            "created": created,
            "modified": modified,
            "author": "admin@trisa.io",
            "editor": Faker().email(),
            "text": lorem.sentence(),
        }
    }


def synthesize_secrets(record):
    """
    For a single record, synthesize sensitive fields
    - signature
    - data
    - chain
    - serial_number

    Returns updated version of the record (dict) with synthetic secrets
    """
    secret = secrets.token_urlsafe(684)
    record["identity_certificate"]["signature"] = secret

    data = secrets.token_urlsafe(3328)
    record["identity_certificate"]["data"] = data

    chain = secrets.token_urlsafe(5920)
    record["identity_certificate"]["chain"] = chain

    serial = secrets.token_urlsafe(24)
    record["identity_certificate"]["serial_number"] = serial

    return record


##########################################################################
# VASP Creation Functions
##########################################################################


def make_verified(vasp, idx, template="fixtures/datagen/templates/verified.json"):
    """
    Populate variable fields in a verified record
    uses `fixtures/datagen/templates/verified.json` as template
    """
    with open(template, "r") as f:
        record = json.load(f)

    state = "VERIFIED"
    country = random.choice(COUNTRIES)

    record["id"] = idx
    record["entity"]["name"] = fake_legal_name(vasp)
    record["entity"]["geographic_addresses"] = [fake_address(country)]
    record["entity"]["national_identification"][
        "national_identifier"
    ] = secrets.token_urlsafe(24)
    record["entity"]["national_identification"]["country_of_issue"] = country
    record["entity"]["country_of_registration"] = country
    record["contacts"]["legal"] = make_person(vasp)
    other = random.choice(
        ["administrative", "technical"]
    )  # billing always unverified for demo purposes
    record["contacts"][other] = make_person(vasp)
    record = synthesize_secrets(record)
    common_name = "trisa." + vasp.lower().split()[0]
    record["common_name"] = common_name + ".io"
    record["identity_certificate"]["subject"]["common_name"] = common_name + ".io"
    record["trisa_endpoint"] = common_name + ".io" + ":123"
    record["website"] = "https://" + common_name + ".io"
    record["vasp_categories"] = [random.sample(VASP_CATEGORIES, random.randint(1, 4))]
    record["established_on"] = Faker().date()
    record["trixo"]["primary_national_jurisdiction"] = country
    record["verification_status"] = state
    record["extra"]["audit_log"] = make_log(state=state)
    record["first_listed"] = record["extra"]["audit_log"][0]["timestamp"]
    record["verified_on"] = record["extra"]["audit_log"][-1]["timestamp"]
    record["last_updated"] = record["extra"]["audit_log"][-1]["timestamp"]
    record["extra"]["review_notes"] = make_notes()

    return record


def make_unverified(
    vasp, idx, state="ERRORED", template="fixtures/datagen/templates/no_cert.json"
):
    """
    Make an unverified record according to the `state`;
    this will be used to create synthetic records for all states other than VERIFIED
    """
    with open(template, "r") as f:
        record = json.load(f)

    country = random.choice(COUNTRIES)

    record["id"] = idx
    record["entity"]["name"] = fake_legal_name(vasp)
    record["entity"]["geographic_addresses"] = [fake_address(country)]
    record["entity"]["national_identification"][
        "national_identifier"
    ] = secrets.token_urlsafe(24)
    record["entity"]["national_identification"]["country_of_issue"] = country
    record["entity"]["country_of_registration"] = country
    record["contacts"]["legal"] = make_person(vasp)
    other = random.choice(["billing", "administrative", "technical"])
    record["contacts"][other] = make_person(vasp)
    common_name = "trisa." + vasp.lower().split()[0]
    record["common_name"] = common_name + ".io"
    record["trisa_endpoint"] = common_name + ".io" + ":123"
    record["website"] = "https://" + common_name + ".io"
    record["vasp_categories"] = [random.sample(VASP_CATEGORIES, random.randint(1, 4))]
    record["established_on"] = Faker().date()
    record["trixo"]["primary_national_jurisdiction"] = country
    record["verification_status"] = state
    record["extra"]["audit_log"] = make_log(state=state)
    record["first_listed"] = record["extra"]["audit_log"][0]["timestamp"]
    record["last_updated"] = record["extra"]["audit_log"][-1]["timestamp"]
    record["extra"]["review_notes"] = make_notes()

    return record


def make_submitted(vasp, idx):
    """
    Populate variable fields in a submitted record
    uses `fixtures/datagen/templates/no_cert.json` as template
    """
    return make_unverified(vasp, idx, state="SUBMITTED")


def make_appealed(vasp, idx):
    """
    Populate variable fields in an appealed record
    uses `fixtures/datagen/templates/no_cert.json` as template
    """
    return make_unverified(vasp, idx, state="APPEALED")


def make_errored(vasp, idx):
    """
    Populate variable fields in an errored record
    uses `fixtures/datagen/templates/no_cert.json` as template
    """
    return make_unverified(vasp, idx, state="ERRORED")


def make_rejected(vasp, idx):
    """
    Populate variable fields in an rejected record
    uses `fixtures/datagen/templates/no_cert.json` as templat
    """
    return make_unverified(vasp, idx, state="REJECTED")


def make_pending(vasp, idx):
    """
    Populate variable fields in an pending review record
    uses `fixtures/datagen/templates/no_cert.json` as template
    """
    return make_unverified(vasp, idx, state="PENDING_REVIEW")


def augment_vasps(fake_names=FAKE_VASPS):
    """
    From the dictionary produced by the load_data function, potentially cleaned by the
    clean_vasp function,
     - Generate new records from keys of FAKE_VASPS, using values to set VASP state
     remaining data is random
     - Add review comments to each record
    Returns synthetic records as a single dictionary
    """
    synthetic_vasps = dict()

    for vasp, state in fake_names.items():
        idx = str(uuid.uuid1())
        if state == "VERIFIED":
            synthetic_vasps[idx] = make_verified(vasp, idx)
        elif state == "ERRORED":
            synthetic_vasps[idx] = make_errored(vasp, idx)
        elif state == "PENDING_REVIEW":
            synthetic_vasps[idx] = make_pending(vasp, idx)
        elif state == "APPEALED":
            synthetic_vasps[idx] = make_appealed(vasp, idx)
        elif state == "REJECTED":
            synthetic_vasps[idx] = make_rejected(vasp, idx)
        elif state == "SUBMITTED":
            synthetic_vasps[idx] = make_submitted(vasp, idx)

    return synthetic_vasps


def store(fakes, directory="fixtures/datagen/synthetic"):
    """
    Save `fakes` dictionary to a new directory (create if not exists)
    Each file should be the name of the fakerized uuid
    Return the path
    """
    if not os.path.exists(directory):
        os.makedirs(directory)

    for idx, fake in fakes.items():
        fname = os.path.join(directory, "vasps::" + idx + ".json")
        with open(fname, "w") as outfile:
            json.dump(fake, outfile, indent=4, sort_keys=True)

    return directory


if __name__ == "__main__":

    fakes = augment_vasps()
    store(fakes)
