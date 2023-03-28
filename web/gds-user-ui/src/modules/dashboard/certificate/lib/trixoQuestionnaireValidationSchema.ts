import * as yup from 'yup';

export const trixoQuestionnaireValidationSchema = yup.object().shape({
  trixo: yup.object().shape({
    primary_national_jurisdiction: yup.string(),
    primary_regulator: yup.string(),
    other_jurisdictions: yup.array().of(
      yup.object().shape({
        country: yup.string(),
        regulator_name: yup.string()
      })
    ),
    financial_transfers_permitted: yup.string().oneOf(['no', 'yes', 'partial']).default('no'),
    has_required_regulatory_program: yup.string().oneOf(['no', 'yes', 'partial']).default('no'),
    conducts_customer_kyc: yup.boolean().default(false),
    kyc_threshold: yup
      .number()
      .transform((value, originalValue) => {
        if (originalValue) {
          const v = originalValue?.replace(/^0+/, '');
          return v?.length > 0 ? Number(v) : 0;
        }
        return value;
      })
      .default(0),
    kyc_threshold_currency: yup.string(),
    must_comply_travel_rule: yup.boolean(),
    applicable_regulations: yup
      .array()
      .of(yup.string())
      .transform((value, originalValue) => {
        if (originalValue) {
          return originalValue.filter((item: any) => item.length > 0);
        }
        return value;
      }),
    compliance_threshold: yup
      .number()
      .transform((value, originalValue: any) => {
        if (originalValue) {
          const v = originalValue?.replace(/^0+/, '');
          return v?.length > 0 ? Number(v) : 0;
        }
        return value;
      })
      .default(0),
    compliance_threshold_currency: yup.string(),
    must_safeguard_pii: yup.boolean().default(false),
    safeguards_pii: yup.boolean().default(false)
  })
});
