import { setupI18n } from '@lingui/core';
import { t } from '@lingui/macro';

import * as yup from 'yup';

const _i18n = setupI18n();

const trisaEndpointPattern = /^([a-zA-Z0-9.-]+):((?!(0))[0-9]+)$/;
const commonNameRegex =
  /^([a-z0-9]+([-a-z0-9]*[a-z0-9]+)?\.){0,}([a-z0-9]+([-a-z0-9]*[a-z0-9]+)?){1,63}(\.[a-z0-9]{2,7})+$/;

export const trisaImplementationValidationSchema = yup.object().shape({
  testnet: yup.object().shape({
    endpoint: yup.string().matches(trisaEndpointPattern, _i18n._(t`TRISA endpoint is not valid`)),
    common_name: yup
      .string()
      .matches(
        commonNameRegex,
        _i18n._(
          t`Common name should not contain special characters, no spaces and must have a dot(.) in it and should have at least 2 characters after the periods`
        )
      )
  }),
  mainnet: yup.object().shape({
    endpoint: yup
      .string()
      .test(
        'uniqueMainetEndpoint',
        _i18n._(t`TestNet and MainNet endpoints should not be the same`),
        (value, ctx: any): any => {
          return ctx.from[1].value.testnet.endpoint !== value;
        }
      )
      .matches(trisaEndpointPattern, _i18n._(t`TRISA endpoint is not valid`)),
    common_name: yup
      .string()
      .matches(
        commonNameRegex,
        _i18n._(
          t`Common name should not contain special characters, no spaces and must have a dot(.) in it and should have at least 2 characters after the periods`
        )
      )
  })
});
