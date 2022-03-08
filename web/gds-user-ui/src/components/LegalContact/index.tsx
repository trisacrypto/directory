import { Heading, Text } from "@chakra-ui/react";
import InputFormControl from "components/ui/InputFormControl";
import PhoneNumberInput from "components/ui/PhoneNumberInput";
import FormLayout from "layouts/FormLayout";

const LegalContact: React.FC = () => {
  return (
    <FormLayout>
      <Heading size="md">Legal/ Compliance Contact (required)</Heading>
      <Text fontStyle="italic">
        Compliance officer or legal contact for requests about the compliance
        requirements and legal status of your organization.
      </Text>
      <InputFormControl
        label="Full Name"
        formHelperText="Preferred name for email communication."
        controlId="fullName"
      />

      <InputFormControl
        label="Email Address"
        formHelperText="Please use the email address associated with your organization."
        controlId="fullName"
        type="email"
      />

      <PhoneNumberInput
        onChange={() => {}}
        label="Phone Number (optional)"
        formHelperText="If supplied, use full phone number with country code."
        controlId="phoneNumber"
      />
    </FormLayout>
  );
};

export default LegalContact;
