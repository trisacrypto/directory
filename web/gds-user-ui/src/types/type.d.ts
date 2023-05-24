import { isoCountries } from '../utils/country';

type StepStatus = 'complete' | 'progress' | 'incomplete';
type NetworkType = 'testnet' | 'mainnet';
type StepperType = 'basic' | 'contacts' | 'legal' | 'trisa' | 'trixo';
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

type Certificate = {
  serial_number: string;
  issued_at: string | Date;
  expires_at: string | Date;
  revoked: boolean;
  details: {
    chain: string;
    data: string;
    endpoint?: string;
    issuer: {
      common_name: string;
      country: string[] | string;
      locality: string[] | string;
      organization: string[] | string;
      organizational_unit: string[] | string;
      postal_code: string[] | string;
      province: string[] | string;
      serial_number: string[] | string;
      street_address: string[] | string;
    };
    not_after: Date | string;
    not_before: Date | string;
    public_key_algorithm: string;
    revoked: boolean;
    serial_number: string;
    signature: string;
    signature_algorithm: string;
    subject: {
      common_name: string;
      country: string[] | string;
      locality: string[] | string;
      organization: string[] | string;
      organizational_unit: string[] | string;
      postal_code: string[] | string;
      province: string[] | string;
      serial_number: string;
      street_address: string[] | string;
    };
    version: string;
  };
};

declare global {
  interface Navigator {
    msSaveBlob?: (blob: any, defaultName?: string) => boolean;
  }
}
