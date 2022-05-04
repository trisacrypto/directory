import { Heading, Text } from '@chakra-ui/react';
import InputFormControl from 'components/ui/InputFormControl';
import PhoneNumberInput from 'components/ui/PhoneNumberInput';
import FormLayout from 'layouts/FormLayout';
import { Controller, useFormContext } from 'react-hook-form';
import get from 'lodash/get';

type ContactFormProps = {
  title: string;
  description: string;
  name: string;
};

const ContactForm: React.FC<ContactFormProps> = ({ title, description, name }) => {
  const { register, control, formState } = useFormContext();
  const { errors } = formState;

  return (
    <FormLayout>
      <Heading size="md">{title}</Heading>
      <Text fontStyle="italic">{description}</Text>
      <InputFormControl
        label="Full Name"
        formHelperText="Preferred name for email communication."
        controlId="fullName"
        isInvalid={get(errors, `${name}.name`)}
        {...register(`${name}.name`)}
      />

      <InputFormControl
        label="Email Address"
        formHelperText={
          get(errors, `${name}.email`)
            ? get(errors, `${name}.email`).message
            : 'Please use the email address associated with your organization.'
        }
        controlId="fullName"
        type="email"
        isInvalid={get(errors, `${name}.email`)}
        {...register(`${name}.email`)}
      />

      <Controller
        control={control}
        name={`${name}.phone`}
        render={({ field: { onChange, value, ref, name: inputName } }) => {
          return (
            <PhoneNumberInput
              ref={ref}
              onChange={(val) => onChange(val)}
              value={value}
              name={inputName}
              isInvalid={get(errors, `${name}.phone`)}
              label="Phone Number (optional)"
              formHelperText="If supplied, use full phone number with country code."
              controlId="phoneNumber"
            />
          );
        }}
      />
    </FormLayout>
  );
};

export default ContactForm;
