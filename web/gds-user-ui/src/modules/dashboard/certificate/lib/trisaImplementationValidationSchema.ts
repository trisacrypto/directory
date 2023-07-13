import { setupI18n } from '@lingui/core';
import { t } from '@lingui/macro';

import * as yup from 'yup';
import {
  trisaImplementationMainnetFieldName,
  trisaImplementationTestnetFieldName
} from './fieldNamesPerSteps';

const _i18n = setupI18n();

const trisaEndpointPattern = /^$|^([a-zA-Z0-9.-]+):((?!(0))[0-9]+)$/;
const commonNameRegex = /^[a-zA-Z0-9]+\.[a-zA-Z0-9]{2,}$/;

export const trisaImplementationValidationSchema = yup.object().shape({
  testnet: yup.object().shape({
    endpoint: yup.string().matches(trisaEndpointPattern, {
      message: _i18n._(t`TRISA endpoint is not valid.`)
    }),
    common_name: yup.string().matches(commonNameRegex, {
      message: _i18n._(
        t`Common name should not contain special characters, no spaces and must have a dot(.) in it and should have at least 2 characters after the periods.`
      )
    })
  }),
  mainnet: yup.object().shape({
    endpoint: yup
      .string()
      .test(
        'uniqueMainetEndpoint',
        _i18n._(t`TestNet and MainNet endpoints should not be the same.`),
        (value, ctx: any) => {
          if (!value) {
            return true;
          }
          return ctx.from[1].value.testnet.endpoint !== value;
        }
      )
      .matches(trisaEndpointPattern, {
        message: _i18n._(t`TRISA endpoint is not valid.`)
      }),
    common_name: yup.string().matches(commonNameRegex, {
      message: _i18n._(
        t`Common name should not contain special characters, no spaces and must have a dot(.) in it and should have at least 2 characters after the periods.`
      )
    })
  }),
  // this field is removed from value object when trisa form is unmounted
  tempField: yup
    .string()

    // should show error when on of mainnet fields is empty
    .when(trisaImplementationMainnetFieldName, {
      is: (...values: string[]) => (values[0] && !values[1]) || (!values[0] && values[1]),
      then: yup.string().required('Mainnet: Common Name or Endpoint is empty.')
    })

    // should show error when on of mainnet fields is empty
    .when(trisaImplementationTestnetFieldName, {
      is: (...values: string[]) => (values[0] && !values[1]) || (!values[0] && values[1]),
      then: yup.string().required('Testnet: Common Name or Endpoint is empty.')
    })
});
