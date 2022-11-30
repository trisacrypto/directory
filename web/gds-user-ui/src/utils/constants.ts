// constants link

export const APP_PATH = {
  HOME: '/',
  GUIDE: '/guide',
  LOGIN: '/auth/login',
  REGISTER: '/auth/register',
  RESET_PASSWORD: '/auth/reset-password',
  PROFILE: '/dashboard/profile',
  PROFILE_EDIT: '/profile/edit',
  DASHBOARD: '/dashboard/overview',
  DASH_CERTIFICATE_REGISTRATION: '/dashboard/certificate/registration',
  CERTIFICATE_REGISTRATION: '/certificate/registration',
  CERTIFICATE_INVENTORY: '/dashboard/certificate/inventory',
  SWITCH_ORGANIZATION: '/dashboard/organization/switch',

};

export const STEPPER_NETWORK = {
  MAINNET: 'mainnet',
  TESTNET: 'testnet',
};

export const APP_STATUS_CODE = {
  NO_CONTENT: 204,
  BAD_REQUEST: 400,
  UNAUTHORIZED: 401,
  FORBIDDEN: 403,
  NOT_FOUND: 404,
  INTERNAL_SERVER_ERROR: 500,
  SERVICE_UNAVAILABLE: 503,
  OK: 200,
};
