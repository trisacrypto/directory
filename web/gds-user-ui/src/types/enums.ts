// enum verification state value
export const enum SubmissionStatus {
  VERIFIED = 'VERIFIED',
  SUBMITTED = 'SUBMITTED',
  NOT_VERIFIED = 'NO_VERIFICATION'
}
// user permission enum value
export const enum USER_PERMISSION {
  READ_COLLABORATOR = 'read:collaborators',
  CREATE_COLLABORATOR = 'create:collaborators',
  UPDATE_COLLABORATOR = 'update:collaborators',
  APPROVE_COLLABORATOR = 'approve:collaborators',
  READ_CERTIFICATE = 'read:certificates',
  CREATE_CERTIFICATE = 'create:certificates',
  UPDATE_CERTIFICATE = 'update:certificates',
  REVOKE_CERTIFICATE = 'revoke:certificates',
  READ_VASP = 'read:vasp',
  CREATE_VASP = 'create:vasp',
  UPDATE_VASP = 'update:vasp',
  CREATE_ORGANIZATIONS = 'create:organizations'
}

export const enum COLLABORATOR_STATUS {
  Pending = 'Pending',
  Confirmed = 'Confirmed'
}

export const enum NetworkType {
  MAINNET = 'mainnet',
  TESTNET = 'testnet'
}

export const enum StepEnum {
  BASIC = 'basic',
  LEGAL = 'legal',
  CONTACTS = 'contacts',
  TRISA = 'trisa',
  TRIXO = 'trixo',
  ALL = 'all'
}
