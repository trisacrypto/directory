import FormLayout from "layouts/FormLayout";
import "react-phone-number-input/style.css";
import PhoneInput, { Value as E164Number } from "react-phone-number-input";
import { Heading, Input } from "@chakra-ui/react";
import PhoneNumberInput from "components/ui/PhoneNumberInput";

type CustomerNumberProps = {};

const CustomerNumber: React.FC<CustomerNumberProps> = () => {
  const handlePhoneInputChange = (value: E164Number) => {
    console.log("[Phone Input] value", value);
  };

  return (
    <FormLayout>
      <Heading size="md">Customer Number</Heading>
      <PhoneNumberInput
        limitMaxLength={true}
        defaultCountry="US"
        controlId="customer_number"
        formHelperText="TRISA specific identity number (UUID), only supplied if you're updating an existing registration request"
        onChange={handlePhoneInputChange}
      />
    </FormLayout>
  );
};

export default CustomerNumber;
