export * from './layout';
export * from './dashboard'

export const ENVIRONMENT = {
    DEV: "development",
    PROD: "production"
}

export const AddressTypeHeaders = {
    ADDRESS_TYPE_CODE_MISC: "Unspecified",
    ADDRESS_TYPE_CODE_HOME: "Residential",
    ADDRESS_TYPE_CODE_BIZZ: "Business",
    ADDRESS_TYPE_CODE_GEOG: "Geographic"
}

export const DIRECTORY = {
    TRISATEST: "TestNet Admin",
    VASP_DIRECTORY: "Production Admin"
}

export const VERIFIED_CONTACT_STATUS = {
    VERIFIED: 'VERIFIED',
    ALTERNATE_VERIFIED: 'ALTERNATE_VERIFIED',
    UNVERIFIED: 'UNVERIFIED'
}

export const VERIFIED_CONTACT_STATUS_LABEL = {
    VERIFIED: 'Verified',
    ALTERNATE_VERIFIED: 'Alternate verified',
    UNVERIFIED: ''
}

export const AUTH_SESSION_KEY = '__SESSION_TOKEN__'

export const LEGAL_PERSON_NAME_IDENTIFIER_TYPE = ['Unspecified', 'Legal Name', 'Short Name', 'Trading Name']

export const validEmailPattern = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/
