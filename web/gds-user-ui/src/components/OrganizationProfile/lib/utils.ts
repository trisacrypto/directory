export const StatusCheck = (status: string) => {
  switch (status) {
    case 'NO_VERIFICATION':
      return 'Not Verified';
    case 'PENDING_REVIEW':
      return 'Pending Review';
    case 'VERIFIED':
      return 'Verified';
    case 'REJECTED':
      return 'Rejected';
    default:
      return 'Pending Registration';
  }
};
