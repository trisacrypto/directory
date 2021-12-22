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
import sys
import json
import uuid
import random
import shutil
import secrets
import tarfile
import datetime

import lorem
from faker import Faker

##########################################################################
# Global Variables
##########################################################################

OUTPUT_DIRECTORY = os.path.join("fixtures", "datagen", "synthetic")
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
    "Echo Funds": "REVIEWED",
    "Foxtrot LLC": "ISSUING_CERTIFICATE",
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
    "Quebec": "READY_TO_SUBMIT",
    "Sierra": "PROCESSING",
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

VASP_CERT_RELATIONSHIPS = {
    "Echo Funds": ["Quebec"],
    "Foxtrot LLC": ["Sierra"],
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

def make_uuid(rng):
    """
    Generate a random UUID using the given random generator.
    """
    return str(uuid.UUID(bytes=bytes(rng.getrandbits(8) for _ in range(16)), version=4))

def make_person(vasp, verified=True, token="", rng=random.Random()):
    """
    Given a string representing a VASP's name, return a valid
    dictionary for a representative of that VASP including faked name
    and contact information.
    """
    fake = Faker()
    name = fake.name()
    domain = rng.choice(DOMAINS)
    email = name.lower().split()[0] + "@" + vasp.replace(" ", "").lower() + domain
    return {
        "name": name,
        "email": email,
        "phone": fake.phone_number(),
        "extra": {
            "@type": "type.googleapis.com/gds.models.v1.GDSContactExtraData",
            "verified": verified,
            "token": token,
            "email_log": [],
        },
    }


def make_dates(first="2021-06-15T05:11:13Z", last="2021-10-25T17:15:43Z", count=3, rng=random.Random()):
    """
    Make `count` number of sequential dates between `first` and `last`,
    return the dates as a list.
    """
    format = "%Y-%m-%dT%H:%M:%SZ"
    start = datetime.datetime.strptime(first, format)
    end = datetime.datetime.strptime(last, format)
    dates = [rng.random() * (end - start) + start for _ in range(count)]
    return [datetime.datetime.strftime(date, format) for date in sorted(dates)]


def make_vasp_log(state="VERIFIED", rng=random.Random()):
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
    dates = make_dates(count=(len(states) - 1), rng=rng)

    for st, dt in zip(states[1:], dates):
        log = dict()
        log.update(VASP_STATE_CHANGES[st])
        log["timestamp"] = dt
        logs.append(log)

    return logs


def make_notes(rng=random.Random()):
    """
    Make fake review notes. Returns a nested dictionary, since notes are a dict
    not a list.
    """
    idx = make_uuid(rng)
    created, modified = make_dates(count=2, rng=rng)
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

def shorten(name):
    """
    Return a shortened version of a name so it can be used in a file path.
    """
    return name.split(" ")[0].lower()

def store(fakes, kind="vasps", directory=OUTPUT_DIRECTORY):
    """
    Save `fakes` dictionary to a new directory (create if not exists)
    Each file should be the name of the fakerized uuid
    Return the path
    """
    if not os.path.exists(directory):
        os.makedirs(directory)

    for idx, fake in fakes.items():
        fname = os.path.join(directory, kind + "::" + shorten(idx) + ".json")
        with open(fname, "w") as outfile:
            json.dump(fake, outfile, indent=4, sort_keys=True)

    return directory

def replace_fixtures():
    """
    Creates a new fakes.tgz file containing the generated fixtures and replaces the
    existing fakes.tgz in the pkg/gds/testdata directory with the new one.
    """
    with tarfile.open("fakes.tgz", "w:gz") as tar:
        tar.add(OUTPUT_DIRECTORY, arcname="synthetic")
    shutil.move("fakes.tgz", os.path.join("pkg", "gds", "testdata", "fakes.tgz"))

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
    rng_country = random.Random(vasp+"country")
    country = rng_country.choice(COUNTRIES)

    record["id"] = idx
    record["entity"]["name"] = fake_legal_name(vasp)
    record["entity"]["geographic_addresses"] = [fake_address(country)]
    record["entity"]["national_identification"][
        "national_identifier"
    ] = secrets.token_urlsafe(24)
    record["entity"]["national_identification"]["country_of_issue"] = country
    record["entity"]["country_of_registration"] = country
    rng_person = random.Random(vasp+"person")
    record["contacts"]["legal"] = make_person(vasp, token="legal_token", rng=rng_person)
    rng_contact = random.Random(vasp+"contact")
    other = rng_contact.choice(
        ["administrative", "technical"]
    )  # billing always unverified for demo purposes
    record["contacts"][other] = make_person(vasp, token=other+"_token", rng=rng_person)
    record = synthesize_secrets(record)
    common_name = "trisa." + vasp.lower().split()[0]
    record["common_name"] = common_name + ".io"
    record["identity_certificate"]["subject"]["common_name"] = common_name + ".io"
    record["trisa_endpoint"] = common_name + ".io" + ":123"
    record["website"] = "https://" + common_name + ".io"
    rng_cat = random.Random(vasp+"cat")
    record["vasp_categories"] = rng_cat.sample(VASP_CATEGORIES, rng_cat.randint(1, 4))
    record["established_on"] = Faker().date()
    record["trixo"]["primary_national_jurisdiction"] = country
    record["verification_status"] = state
    rng_audit = random.Random(vasp+"audit")
    record["extra"]["audit_log"] = make_vasp_log(state=state, rng=rng_audit)
    record["first_listed"] = record["extra"]["audit_log"][0]["timestamp"]
    record["verified_on"] = record["extra"]["audit_log"][-1]["timestamp"]
    record["last_updated"] = record["extra"]["audit_log"][-1]["timestamp"]
    rng_notes = random.Random(vasp+"notes")
    record["extra"]["review_notes"] = make_notes(rng=rng_notes)

    return record


def make_unverified(
    vasp, idx, state="ERRORED", template="fixtures/datagen/templates/no_cert.json"):
    """
    Make an unverified record according to the `state`;
    this will be used to create synthetic records for all states other than VERIFIED
    """
    with open(template, "r") as f:
        record = json.load(f)

    # A VASP that is not at least EMAIL_VERIFIED cannot have verified contacts
    email_verified = (state != "SUBMITTED")

    rng_country = random.Random(vasp)
    country = rng_country.choice(COUNTRIES)

    record["id"] = idx
    record["entity"]["name"] = fake_legal_name(vasp)
    record["entity"]["geographic_addresses"] = [fake_address(country)]
    record["entity"]["national_identification"][
        "national_identifier"
    ] = secrets.token_urlsafe(24)
    record["entity"]["national_identification"]["country_of_issue"] = country
    record["entity"]["country_of_registration"] = country
    rng_person = random.Random(vasp+"person")
    record["contacts"]["legal"] = make_person(vasp, verified=email_verified, token="legal_token", rng=rng_person)
    rng_contact = random.Random(vasp+"contact")
    other = rng_contact.choice(["billing", "administrative", "technical"])
    record["contacts"][other] = make_person(vasp, verified=email_verified, token=other+"_token", rng=rng_person)
    common_name = "trisa." + vasp.lower().split()[0]
    record["common_name"] = common_name + ".io"
    record["trisa_endpoint"] = common_name + ".io" + ":123"
    record["website"] = "https://" + common_name + ".io"
    rng_cat = random.Random(vasp+"cat")
    record["vasp_categories"] = rng_cat.sample(VASP_CATEGORIES, rng_cat.randint(1, 4))
    record["established_on"] = Faker().date()
    record["trixo"]["primary_national_jurisdiction"] = country
    record["verification_status"] = state
    rng_audit = random.Random(vasp+"audit")
    record["extra"]["audit_log"] = make_vasp_log(state=state, rng=rng_audit)
    record["first_listed"] = record["extra"]["audit_log"][0]["timestamp"]
    record["last_updated"] = record["extra"]["audit_log"][-1]["timestamp"]
    rng_notes = random.Random(vasp+"notes")
    record["extra"]["review_notes"] = make_notes(rng=rng_notes)

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

def make_reviewed(vasp, idx):
    """
    Populate variable fields in a reviewed record
    uses `fixtures/datagen/templates/no_cert.json` as template
    """
    return make_unverified(vasp, idx, state="REVIEWED")

def make_issuing(vasp, idx):
    """
    Populate variable fields in a certificate issued record
    uses `fixtures/datagen/templates/no_cert.json` as template
    """
    return make_unverified(vasp, idx, state="ISSUING_CERTIFICATE")

def augment_vasps(fake_names=FAKE_VASPS):
    """
    Generate new records from keys of FAKE_VASPS, using values to set VASP state
    The remaining data is random. Add review comments to each record
    Returns synthetic records as a single dictionary
    """
    rng = random.Random("vasps")
    synthetic_vasps = dict()

    for vasp, state in fake_names.items():
        idx = make_uuid(rng)
        if state == "VERIFIED":
            synthetic_vasps[vasp] = make_verified(vasp, idx)
        elif state == "ERRORED":
            synthetic_vasps[vasp] = make_errored(vasp, idx)
        elif state == "PENDING_REVIEW":
            synthetic_vasps[vasp] = make_pending(vasp, idx)
        elif state == "REVIEWED":
            synthetic_vasps[vasp] = make_reviewed(vasp, idx)
        elif state == "ISSUING_CERTIFICATE":
            synthetic_vasps[vasp] = make_issuing(vasp, idx)
        elif state == "APPEALED":
            synthetic_vasps[vasp] = make_appealed(vasp, idx)
        elif state == "REJECTED":
            synthetic_vasps[vasp] = make_rejected(vasp, idx)
        elif state == "SUBMITTED":
            synthetic_vasps[vasp] = make_submitted(vasp, idx)
        else:
            print("Skipping unrecognized state: %s", state)

    return synthetic_vasps


##########################################################################
# CertReq Creation Functions
##########################################################################


def make_common_name(cert, rng=random.Random()):
    """
    Make a synthetic but well-structured common name
    """
    return cert + "." + rng.choice(URLWORDS) + rng.choice(DOMAINS)


def make_completed(cert, idx, template="fixtures/datagen/templates/cert_req.json", rng=random.Random()):
    """
    Make a cert req in the COMPLETED state
    """
    with open(template, "r") as f:
        record = json.load(f)

    record["id"] = idx
    record["vasp"] = make_uuid(rng)
    common_name = make_common_name(cert)
    record["common_name"] = common_name
    record["status"] = "COMPLETED"
    batch = str(rng.randint(100000, 999999))
    record["batch_id"] = batch
    network = rng.choice(NETWORKS)
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


def make_initialized(cert, idx, template="fixtures/datagen/templates/cert_req.json", rng=random.Random()):
    """
    Make a cert req in the INITIALIZED state
    """
    with open(template, "r") as f:
        record = json.load(f)

    record["id"] = idx
    record["vasp"] = make_uuid(rng)
    common_name = make_common_name(cert)
    record["common_name"] = common_name
    record["status"] = "INITIALIZED"
    start = make_dates(count=1)[0]
    record["created"] = start
    record["modified"] = start
    record["audit_log"] = make_cert_log("INITIALIZED", start, start)
    return record

def make_ready_to_submit(cert, idx, template="fixtures/datagen/templates/cert_req.json", rng=random.Random()):
    """
    Make a cert req in the READY_TO_SUBMIT state
    """
    with open(template, "r") as f:
        record = json.load(f)

    record["id"] = idx
    record["vasp"] = make_uuid(rng)
    common_name = make_common_name(cert)
    record["common_name"] = common_name
    record["status"] = "READY_TO_SUBMIT"
    start, end = make_dates(count=2)
    record["created"] = start
    record["modified"] = end
    record["audit_log"] = make_cert_log("READY_TO_SUBMIT", start, end)
    return record

def make_processing(cert, idx, template="fixtures/datagen/templates/cert_req.json", rng=random.Random()):
    """
    Make a cert req in the PROCESSING state
    """
    with open(template, "r") as f:
        record = json.load(f)

    record["id"] = idx
    record["vasp"] = make_uuid(rng)
    common_name = make_common_name(cert)
    record["common_name"] = common_name
    record["status"] = "PROCESSING"
    batch = str(rng.randint(100000, 999999))
    record["batch_id"] = batch
    start, end = make_dates(count=2)
    record["created"] = start
    record["modified"] = end
    record["audit_log"] = make_cert_log("PROCESSING", start, end)
    return record

def make_cr_errored(cert, idx, template="fixtures/datagen/templates/cert_req.json", rng=random.Random()):
    """
    Make a cert req in the CR_ERRORED state
    """
    with open(template, "r") as f:
        record = json.load(f)

    record["id"] = idx
    record["vasp"] = make_uuid(rng)
    common_name = make_common_name(cert)
    record["common_name"] = common_name
    record["status"] = "CR_ERRORED"
    batch = str(rng.randint(100000, 999999))
    record["batch_id"] = batch
    network = rng.choice(NETWORKS)
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


def make_cr_rejected(cert, idx, template="fixtures/datagen/templates/cert_req.json", rng=random.Random()):
    """
    Make a cert req in the CR_REJECTED state
    """
    with open(template, "r") as f:
        record = json.load(f)

    record["id"] = idx
    record["vasp"] = make_uuid(rng)
    common_name = make_common_name(cert)
    record["common_name"] = common_name
    record["status"] = "CR_REJECTED"
    batch = str(rng.randint(100000, 999999))
    record["batch_id"] = batch
    network = rng.choice(NETWORKS)
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
    rng = random.Random("certs")
    synthetic_certs = dict()

    for cert, state in fake_names.items():
        idx = make_uuid(rng)
        cert = cert.lower()
        if state == "INITIALIZED":
            synthetic_certs[cert] = make_initialized(cert, idx)
        elif state == "READY_TO_SUBMIT":
            synthetic_certs[cert] = make_ready_to_submit(cert, idx)
        elif state == "PROCESSING":
            synthetic_certs[cert] = make_processing(cert, idx)
        elif state == "COMPLETED":
            synthetic_certs[cert] = make_completed(cert, idx)
        elif state == "CR_ERRORED":
            synthetic_certs[cert] = make_cr_errored(cert, idx)
        elif state == "CR_REJECTED":
            synthetic_certs[cert] = make_cr_rejected(cert, idx)
        else:
            print("Skipping unrecognized state: %s", state)

    return synthetic_certs

def add_vasp_cert_relationships(vasps, certs):
    """
    Add predefined relationships between VASPs and certificate requests.
    """
    for vasp_name, cert_names in VASP_CERT_RELATIONSHIPS.items():
        cert_ids = []
        for c in cert_names:
            name = c.lower()
            cert_ids.append(certs[name]["id"])
            certs[name]["vasp"] = vasps[vasp_name]["id"]
            certs[name]["common_name"] = vasps[vasp_name]["common_name"]
        vasps[vasp_name]["extra"]["certificate_requests"] = cert_ids

if __name__ == "__main__":
    replace = False
    if len(sys.argv) > 1:
        if sys.argv[1] == "--help":
            print("Usage: python fixtures/datagen/fakerize.py [--help|--replace]")
            print("  --help: print this message")
            print("  --replace: generate and replace existing fixtures in pkg/gds/testdata")
            sys.exit(0)
        elif sys.argv[1] == "--replace":
            replace = True
        else:
            print("Unknown argument: %s", sys.argv[1])
            sys.exit(1)

    if os.path.exists(OUTPUT_DIRECTORY):
        shutil.rmtree(OUTPUT_DIRECTORY)
        
    fake_vasps = augment_vasps()
    fake_certs = augment_certs()
    add_vasp_cert_relationships(fake_vasps, fake_certs)

    store(fake_vasps, kind="vasps")
    store(fake_certs, kind="certreqs")

    if replace:
        replace_fixtures()
        print("Successfully replaced pkg/gds/testdata/fakes.tgz")
