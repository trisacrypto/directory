import { Heading } from "@chakra-ui/react";
import InputFormControl from "components/ui/InputFormControl";
import FormLayout from "layouts/FormLayout";

type CountryOfRegistrationProps = {};
const CountryOfRegistration: React.FC<CountryOfRegistrationProps> = () => {
  return (
    <FormLayout>
      <Heading size="md">Country of Registration</Heading>
      <InputFormControl controlId="country_of_registration" />
    </FormLayout>
  );
};

export default CountryOfRegistration;
