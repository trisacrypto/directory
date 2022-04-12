import _ from 'lodash';
export const findStepKey = (steps: any, key: number) =>
  steps.filter((step: any) => step.key === key);

export const isValidUuid = (str: string) => {
  // Regular expression to check if string is a valid UUID
  const regexExp =
    /^[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}$/gi;
  return regexExp.test(str);
};

export const getStepData = (steps: any, key: number): TStep | undefined => {
  const s = findStepKey(steps, key);
  if (s && s?.length === 1) {
    return s[0].data;
  }
  return undefined;
};

export const hasStepError = (steps: any): boolean => {
  const s = steps.filter((step: any) => step.status === 'error');
  return s.length > 0;
};

export const getStepDatas = (steps: any) => {
  const s = steps
    ?.map((step: any) => step.data)
    .reduce((acc: any, cur: any) => ({ ...acc, ...cur }), {});

  return { ...s };
};

export const mapStepsDataToFormRequest = (steps: any) => {
  const s = getStepDatas(steps);

  return {
    entity: {
      name: {
        name_identifiers: s['entity.name.name_identifiers'].map((nameIdentifier: any) => ({
          legal_person_name: nameIdentifier.name_identifier,
          legal_person_name_identifier_type: nameIdentifier.name_identifier_type
        })),
        local_name_identifiers: s['entity.name.local_name_identifiers'].map(
          (nameIdentifier: any) => ({
            local_name_identifier: nameIdentifier.name_identifier,
            local_name_identifier_type: nameIdentifier.name_identifier_type
          })
        ),
        phonetic_name_identifiers: s['entity.name.phonetic_name_identifiers'].map(
          (nameIdentifier: any) => ({
            phonetic_name_identifier: nameIdentifier.name_identifier,
            phonetic_name_identifier_type: nameIdentifier.name_identifier_type
          })
        )
      },
      geographic_adressess: s['entity.geographic_adressess'].map((geographicAddress: any) => ({
        address_line_1: geographicAddress.address_line_1,
        address_line_2: geographicAddress.address_line_2,
        address_line_3: geographicAddress.address_line_3
      }))
    },
    contacts: {
      administrative: {
        email: s['contacts.administrative.email'],
        phone: s['contacts.administrative.phone'],
        name: s['contacts.administrative.name']
      },
      technical: {
        email: s['contacts.technical.email'],
        phone: s['contacts.technical.phone'],
        name: s['contacts.technical.name']
      },
      billing: s['contacts.billing']
        ? {
            email: s['contacts.billing.email'],
            phone: s['contacts.billing.phone'],
            name: s['contacts.billing.name']
          }
        : {},
      legal: s['contacts.legal.email']
        ? {
            email: s['contacts.legal.email'],
            phone: s['contacts.legal.phone'],
            name: s['contacts.legal.name']
          }
        : {}
    },
    webiste: s.website,
    business_category: s.business_category,
    vasp_categories: s.vasp_categories,
    established_on: s.established_on,
    trixo: {
      primary_national_juridisction: s['trixo.primary_national_juridisction'],
      primary_regulator: s['trixo.primary_regulator'],
      other_jurisdictions: s['trixo.other_jurisdictions'].map((otherJurisdiction: any) => ({
        country: otherJurisdiction.country,
        regulator_name: otherJurisdiction.regulator_name
      })),
      financial_transfers_permitted: s['trixo.financial_transfers_permitted'],
      has_required_regulatory_program: s['trixo.has_required_regulatory_program'],
      conducts_customer_kyc: s['trixo.conducts_customer_kyc'],
      has_required_compliance_program: s['trixo.has_required_compliance_program']
    }
  };
};

export const getValueByPathname = (obj: Record<string, any>, path: string) => {
  return _.get(obj, path);
};

export const getDomain = (url: string | URL) => {
  try {
    const _url = new URL(url);
    return _url.hostname.replace('www.', '');
  } catch (error) {
    console.error('[error]', error);
    throw new Error('Invalid URL format');
  }
};
