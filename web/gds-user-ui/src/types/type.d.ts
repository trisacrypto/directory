import { isoCountries } from './../utils/country';
// define common type
type TStep = {
  status: StepStatus;
  key?: number;
  data?: any;
};

type StepStatus = 'complete' | 'progress' | 'incomplete';

type IsoCountryCode = keyof typeof isoCountries;

type RegistrationAuthority = {
  option: string;
  country: string;
  register: string;
  organization: string;
  website: string;
  comments: string;
  isDisabled?: boolean;
};
