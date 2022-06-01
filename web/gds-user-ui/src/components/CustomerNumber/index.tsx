import FormLayout from 'layouts/FormLayout';
import 'react-phone-number-input/style.css';
import PhoneInput, { Value as E164Number } from 'react-phone-number-input';
import { Heading, Input } from '@chakra-ui/react';
import PhoneNumberInput from 'components/ui/PhoneNumberInput';
import { t } from '@lingui/macro';
import { Trans } from '@lingui/react';

type CustomerNumberProps = {};

const CustomerNumber: React.FC<CustomerNumberProps> = () => {
  const handlePhoneInputChange = (value: E164Number) => {};

  return (
    <FormLayout>
      <Heading size="md">
        <Trans id="Customer Number">Customer Number</Trans>
      </Heading>
      <PhoneNumberInput
        limitMaxLength={true}
        defaultCountry="US"
        controlId="customer_number"
        formHelperText={t`TRISA specific identity number (UUID), only supplied if you're updating an existing registration request`}
        onChange={handlePhoneInputChange}
      />
    </FormLayout>
  );
};

export default CustomerNumber;
