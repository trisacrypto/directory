import { isoCountries } from './../utils/country';

type StepStatus = 'complete' | 'progress' | 'incomplete';

type IsoCountryCode = keyof typeof isoCountries;

type RegistrationAuthority = {
  option: string;
  country: string;
  register: string;
  organization: string;
  website: string;
  jurisdiction: string;
  country_name: string;
  comments: string;
  isDisabled?: boolean;
};

type Locales = 'en' | 'fr' | 'ja' | 'de' | 'zh';
