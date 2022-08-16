import _ from 'lodash';
import { RegistrationAuthority, StepStatus } from 'types/type';
import registrationAuthority from './registration-authority.json';
import auth0 from 'auth0-js';
import getAuth0Config from 'application/config/auth0';
import * as Sentry from '@sentry/react';
import { getRegistrationDefaultValue } from 'modules/dashboard/registration/utils';
import { clearCookies } from 'utils/cookies';
import { defaultIfEmpty } from 'rxjs';
const DEFAULT_REGISTRATION_AUTHORITY = 'RA777777';
export const findStepKey = (steps: any, key: number) =>
  steps?.filter((step: any) => step.key === key);

export const isValidUuid = (str: string) => {
  // Regular expression to check if string is a valid UUID
  const regexExp =
    /^[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}$/gi;
  return regexExp.test(str);
};

export const getStepStatus = (steps: any, key: number): StepStatus | undefined => {
  const s = findStepKey(steps, key);
  if (s && s?.length === 1) {
    return s[0].status;
  }
  return undefined;
};

export const hasStepError = (steps: any): boolean => {
  const s = steps.filter((step: any) => step.status === 'error');
  return s.length > 0;
};

export const getValueByPathname = (obj: Record<string, any>, path: string) => {
  return _.get(obj, path);
};

export const getDomain = (url: string | URL): string | null => {
  try {
    const _url = new URL(url);
    return _url?.hostname?.replace('www.', '');
  } catch (error) {
    console.error('[error]', error);
    return null;
  }
};

export const getRegistrationAuthorities = () => [...new Set(registrationAuthority)];

export const getRegistrationAuthoritiesOptions = (country?: any) => {
  const newArray = [...Array.from(new Set(registrationAuthority))];
  if (country) {
    return newArray
      .filter(
        (v: RegistrationAuthority) =>
          v.country === country || v.option === DEFAULT_REGISTRATION_AUTHORITY
      )
      .map((v: RegistrationAuthority) => {
        const label = v.organization ? `${v.option} - ${v.organization}` : `${v.option}`;
        const l =
          v.jurisdiction && v.jurisdiction !== v.country_name
            ? `${label} - ${v.jurisdiction}`
            : label;
        return {
          value: v.option,
          label: l,
          isDisabled: !!v.isDisabled
        };
      });
  }
  return newArray.map((v: RegistrationAuthority) => {
    const label = v.organization ? `${v.option} - ${v.organization}` : `${v.option}`;
    return {
      value: v.option,
      label,
      isDisabled: !!v.isDisabled
    };
  });
};

export const mapTrixoFormForBff = (data: any) => {
  const { trixo } = data;
  const { applicable_regulations, other_jurisdictions } = trixo;

  const cleanAppRegulation = applicable_regulations.reduce((acc: any, value: any) => {
    if (value.name.length > 0) {
      acc.push(value.name);
    }
    return acc;
  }, []);
  const cleanOtherJurisdiction = other_jurisdictions.filter(
    (o: any) => o?.country?.length > 0 && o?.regulator_name?.length > 0
  );

  return {
    ...data,
    trixo: {
      ...trixo,
      applicable_regulations: cleanAppRegulation.length > 0 ? cleanAppRegulation : [],
      other_jurisdictions: cleanOtherJurisdiction.length > 0 ? cleanOtherJurisdiction : []
    }
  };
};

export const hasValue = (obj: Record<string, any>): boolean => {
  return obj && Object.values(obj).some(Boolean);
};

export const getColorScheme = (status: string) => {
  if (status === 'yes' || status) {
    return 'cyan';
  } else {
    return '#eee';
  }
};

export function currencyFormatter(
  amount: number | bigint,
  { style = 'currency', currency = 'USD' }: Intl.NumberFormatOptions = {}
) {
  const formatedAmount = new Intl.NumberFormat('en-US', {
    style,
    currency
  });

  return formatedAmount.format(amount);
}

export const getRefreshToken = () => {
  const auth0Config = getAuth0Config();
  const authWeb = new auth0.WebAuth(auth0Config);
  return new Promise((resolve, reject) => {
    authWeb.checkSession(
      {
        scope: 'read:current_user'
      },
      (err: any, authResult: any) => {
        if (err) {
          reject(err);
        } else {
          resolve(authResult);
        }
      }
    );
  });
};
export const handleError = (error: any, customMessage?: string) => {
  Sentry.captureMessage(customMessage || error);
  Sentry.captureException(error);
  console.log(error.response);
  if (error.response.status === 401 || error.response.status === 403) {
    clearCookies();

    switch (error.response.status) {
      case '401':
        window.location.href = `/auth/login?q=token_expired`;
        break;
      case '403':
        window.location.href = `/auth/login?q=unauthorized`;
        break;
      default:
        window.location.href = `/auth/login?error_description=${error.response.data.error}`;
    }

    // if (error.response.status === 401) {
    //   window.location.href = `/auth/login?q=token_expired`;
    // }
    // if (error.response.status === 403) {
    //   window.location.href = `/auth/login?q=unauthorized`;
    // }
  }
};

// uppercased first letter
export const upperCaseFirstLetter = (str: any) => {
  return str?.charAt(0)?.toUpperCase() + str?.slice(1);
};

// load default value for registration
export const loadStepperDefaultValue = () => {
  const fetchDefaultValue = async () => {
    const response = await getRegistrationDefaultValue();
    return response.stepper;
  };
  const defaultValue = fetchDefaultValue();
  return {
    ...defaultValue,
    hasReachSubmitStep: false
  };
};
