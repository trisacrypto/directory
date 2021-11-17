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
import random
import secrets
import datetime

import lorem
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
VASP_STATE_CHANGES = {
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

NETWORKS = ["trisatest.net", "vaspdirectory.net"]
URLWORDS = [
    "cacao",
    "pepper",
    "jackolantern",
    "bones",
    "countably",
    "roosevelt",
    "mountain",
    "ace",
    "lighthouse",
    "tauceti",
    "planetary",
    "colloquial",
    "sculptural",
    "estimator",
    "geodistributed",
    "princeton",
    "excelsior",
    "gormandize",
    "wistful",
    "philosophers",
    "hellenic",
]
FAKE_CERTS = {
    "Papa": "INITIALIZED",
    "Quebec": "COMPLETED",
    "Sierra": "CR_REJECTED",
    "Tango": "CR_ERRORED",
    "Uniform": "COMPLETED",
    "Victor": "COMPLETED",
    "Whiskey": "CR_REJECTED",
    "XRay": "INITIALIZED",
    "Yankee": "INITIALIZED",
    "Zulu": "COMPLETED",
}

CERT_STATE_CHANGES = {
    "INITIALIZED": {
        "previous_state": "INITIALIZED",
        "current_state": "INITIALIZED",
        "description": "created certificate request",
        "source": "automated",
    },
    "READY_TO_SUBMIT": {
        "previous_state": "INITIALIZED",
        "current_state": "READY_TO_SUBMIT",
        "description": "registration request received",
        "source": "admin@rotational.io",
    },
    "PROCESSING": {
        "previous_state": "READY_TO_SUBMIT",
        "current_state": "PROCESSING",
        "description": "certificate submitted",
        "source": "automated",
    },
    "CR_REJECTED": {
        "previous_state": "PROCESSING",
        "current_state": "CR_REJECTED",
        "description": "failed domain verification",
        "source": "automated",
    },
    "CR_ERRORED": {
        "previous_state": "PROCESSING",
        "current_state": "CR_ERRORED",
        "description": "certificate limit reached",
        "source": "automated",
    },
    "DOWNLOADING": {
        "previous_state": "PROCESSING",
        "current_state": "DOWNLOADING",
        "description": "certificate ready for download",
        "source": "automated",
    },
    "DOWNLOADED": {
        "previous_state": "DOWNLOADING",
        "current_state": "DOWNLOADED",
        "description": "certificate downloaded",
        "source": "automated",
    },
    "COMPLETED": {
        "previous_state": "DOWNLOADED",
        "current_state": "COMPLETED",
        "description": "certificate request complete",
        "source": "automated",
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


def make_vasp_log(state="VERIFIED"):
    """
    Create a fake audit log depending on the dates provided
    and the setting of `state`.
    Returns a list of dictionaries.
    """
    logs = []

    states = [state]
    current_state = state
    prior_state = None
    while current_state in VASP_STATE_CHANGES:
        prior_state = VASP_STATE_CHANGES[current_state]["previous_state"]
        states.insert(0, prior_state)
        current_state = prior_state

    # skip the first data and state, since it's NO_VERIFICATION and doesn't get a log
    dates = make_dates(count=(len(states) - 1))

    for st, dt in zip(states[1:], dates):
        log = dict()
        log.update(VASP_STATE_CHANGES[st])
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


def store(fakes, kind="vasps", directory="fixtures/datagen/synthetic"):
    """
    Save `fakes` dictionary to a new directory (create if not exists)
    Each file should be the name of the fakerized uuid
    Return the path
    """
    if not os.path.exists(directory):
        os.makedirs(directory)

    for idx, fake in fakes.items():
        fname = os.path.join(directory, kind + "::" + idx + ".json")
        with open(fname, "w") as outfile:
            json.dump(fake, outfile, indent=4, sort_keys=True)

    return directory


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
    record["vasp_categories"] = random.sample(VASP_CATEGORIES, random.randint(1, 4))
    record["established_on"] = Faker().date()
    record["trixo"]["primary_national_jurisdiction"] = country
    record["verification_status"] = state
    record["extra"]["audit_log"] = make_vasp_log(state=state)
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
    record["vasp_categories"] = random.sample(VASP_CATEGORIES, random.randint(1, 4))
    record["established_on"] = Faker().date()
    record["trixo"]["primary_national_jurisdiction"] = country
    record["verification_status"] = state
    record["extra"]["audit_log"] = make_vasp_log(state=state)
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
    Generate new records from keys of FAKE_VASPS, using values to set VASP state
    The remaining data is random. Add review comments to each record
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
        else:
            print("Skipping unrecognized state: %s", state)

    return synthetic_vasps


##########################################################################
# CertReq Creation Functions
##########################################################################


def make_common_name(cert):
    """
    Make a synthetic but well-structured common name
    """
    return cert + "." + random.choice(URLWORDS) + random.choice(DOMAINS)


def make_completed(cert, idx, template="fixtures/datagen/templates/cert_req.json"):
    """
    Make a cert req in the COMPLETED state
    """
    with open(template, "r") as f:
        record = json.load(f)

    record["id"] = idx
    record["vasp"] = str(uuid.uuid1())
    common_name = make_common_name(cert)
    record["common_name"] = common_name
    record["status"] = "COMPLETED"
    batch = str(random.randint(100000, 999999))
    record["batch_id"] = batch
    network = random.choice(NETWORKS)
    record[
        "batch_name"
    ] = f"{network} certificate request for {common_name} (id: {idx})"
    record["batch_status"] = "READY_FOR_DOWNLOAD"
    record["order_number"] = batch
    start, end = make_dates(count=2)
    record["creation_date"] = end
    record["created"] = start
    record["modified"] = end
    record["audit_log"] = make_cert_log("COMPLETED", start, end)
    return record


def make_initialized(cert, idx, template="fixtures/datagen/templates/cert_req.json"):
    """
    Make a cert req in the INITIALIZED state
    """
    with open(template, "r") as f:
        record = json.load(f)

    record["id"] = idx
    record["vasp"] = str(uuid.uuid1())
    common_name = make_common_name(cert)
    record["common_name"] = common_name
    record["status"] = "INITIALIZED"
    start = make_dates(count=1)[0]
    record["created"] = start
    record["modified"] = start
    record["audit_log"] = make_cert_log("INITIALIZED", start, start)
    return record


def make_cr_errored(cert, idx, template="fixtures/datagen/templates/cert_req.json"):
    """
    Make a cert req in the CR_ERRORED state
    """
    with open(template, "r") as f:
        record = json.load(f)

    record["id"] = idx
    record["vasp"] = str(uuid.uuid1())
    common_name = make_common_name(cert)
    record["common_name"] = common_name
    record["status"] = "CR_ERRORED"
    batch = str(random.randint(100000, 999999))
    record["batch_id"] = batch
    network = random.choice(NETWORKS)
    record[
        "batch_name"
    ] = f"{network} certificate request for {common_name} (id: {idx})"
    record["batch_status"] = "FAILED"
    record["order_number"] = batch
    start, end = make_dates(count=2)
    record["created"] = start
    record["modified"] = start
    record["audit_log"] = make_cert_log("CR_ERRORED", start, end)
    return record


def make_cr_rejected(cert, idx, template="fixtures/datagen/templates/cert_req.json"):
    """
    Make a cert req in the CR_REJECTED state
    """
    with open(template, "r") as f:
        record = json.load(f)

    record["id"] = idx
    record["vasp"] = str(uuid.uuid1())
    common_name = make_common_name(cert)
    record["common_name"] = common_name
    record["status"] = "CR_REJECTED"
    batch = str(random.randint(100000, 999999))
    record["batch_id"] = batch
    network = random.choice(NETWORKS)
    record[
        "batch_name"
    ] = f"{network} certificate request for {common_name} (id: {idx})"
    record["batch_status"] = "REJECTED"
    record["order_number"] = batch
    start, end = make_dates(count=2)
    record["created"] = start
    record["modified"] = start
    record["reject_reason"] = lorem.sentence()
    record["audit_log"] = make_cert_log("CR_REJECTED", start, end)
    return record


def make_cert_log(state, start, end):
    """
    Return a list of dictionaries representing a synthetic but plausible audit log
    """
    logs = []

    states = [state]
    current_state = state
    prior_state = None
    while current_state in CERT_STATE_CHANGES:
        prior_state = CERT_STATE_CHANGES[current_state]["previous_state"]
        states.insert(0, prior_state)
        if prior_state == "INITIALIZED":
            current_state = "STOP"
        else:
            current_state = prior_state

    dates = make_dates(first=start, last=end, count=len(states))

    for st, dt in zip(states, dates):
        log = dict()
        log.update(CERT_STATE_CHANGES[st])
        log["timestamp"] = dt
        logs.append(log)

    return logs


def augment_certs(fake_names=FAKE_CERTS):
    """
    Generate new records from keys of FAKE_CERTS, using values to set cert state
    The remaining data is random. Add audit logs to each record
    Returns synthetic records as a single dictionary
    """
    synthetic_certs = dict()

    for cert, state in fake_names.items():
        idx = str(uuid.uuid1())
        cert = cert.lower()
        if state == "INITIALIZED":
            synthetic_certs[idx] = make_initialized(cert, idx)
        elif state == "COMPLETED":
            synthetic_certs[idx] = make_completed(cert, idx)
        elif state == "CR_ERRORED":
            synthetic_certs[idx] = make_cr_errored(cert, idx)
        elif state == "CR_REJECTED":
            synthetic_certs[idx] = make_cr_rejected(cert, idx)
        else:
            print("Skipping unrecognized state: %s", state)

    return synthetic_certs


if __name__ == "__main__":
    fakes = augment_vasps()
    store(fakes, kind="vasps")

    fakes = augment_certs()
    store(fakes, kind="certreqs")
