import { Heading, Text } from "@chakra-ui/react";
import InputFormControl from "components/ui/InputFormControl";
import PhoneNumberInput from "components/ui/PhoneNumberInput";
import FormLayout from "layouts/FormLayout";

type ContactFormProps = {
  title: string;
  description: string;
};

const ContactForm: React.FC<ContactFormProps> = ({ title, description }) => {
  return (
    <FormLayout>
      <Heading size="md">{title}</Heading>
      <Text fontStyle="italic">{description}</Text>
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

export default ContactForm;
