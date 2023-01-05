export const Status = {
  SUBMITTED: 'SUBMITTED',
  EMAIL_VERIFIED: 'EMAIL_VERIFIED',
  PENDING_REVIEW: 'PENDING_REVIEW',
  REVIEWED: 'REVIEWED',
  ISSUING_CERTIFICATE: 'ISSUING_CERTIFICATE',
  VERIFIED: 'VERIFIED',
  REJECTED: 'REJECTED',
  APPEALED: 'APPEALED',
  ERRORED: 'ERRORED',
  NO_VERIFICATION: 'NO_VERIFICATION',
} as const;

export const StatusLabel = {
  SUBMITTED: 'Submitted',
  EMAIL_VERIFIED: 'Email Verified',
  PENDING_REVIEW: 'Pending Review',
  REVIEWED: 'Reviewed',
  ISSUING_CERTIFICATE: 'Issuing Certificate',
  VERIFIED: 'Verified',
  REJECTED: 'Rejected',
  APPEALED: 'Appealed',
  ERRORED: 'Errored',
  NO_VERIFICATION: 'No Verification',
} as const;
