import { Heading, Text } from '@chakra-ui/react';
import InputFormControl from 'components/ui/InputFormControl';
import PhoneNumberInput from 'components/ui/PhoneNumberInput';
import FormLayout from 'layouts/FormLayout';
import { Controller, useFormContext } from 'react-hook-form';
import get from 'lodash/get';
import { t } from '@lingui/macro';

type ContactFormProps = {
  title: string;
  description: string;
  name: string;
};

const ContactForm: React.FC<ContactFormProps> = ({ title, description, name }) => {
  const { register, control, formState } = useFormContext();
  const { errors } = formState;
  const getPhoneMessageHint = () => {
    if (name === 'contacts.legal') {
      return 'A business phone number is required to complete physical verification for MainNet registration. Please provide a phone number where the Legal/ Compliance contact can be contacted.';
    }
    return 'If supplied, use full phone number with country code.';
  };

  return (
    <FormLayout>
      <Heading size="md">{title}</Heading>
      <Text fontStyle="italic">{description}</Text>
      <InputFormControl
        label={t`Full Name`}
        formHelperText={t`Preferred name for email communication.`}
        controlId="fullName"
        isInvalid={get(errors, `${name}.name`)}
        {...register(`${name}.name`)}
      />

      <InputFormControl
        label={t`Email Address`}
        formHelperText={
          get(errors, `${name}.email`)
            ? get(errors, `${name}.email`).message
            : t`Please use the email address associated with your organization.`
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
              label="Phone Number "
              formHelperText={
                get(errors, `${name}.phone`)
                  ? get(errors, `${name}.phone`).message
                  : getPhoneMessageHint()
              }
              controlId="phoneNumber"
            />
          );
        }}
      />
    </FormLayout>
  );
};

export default ContactForm;
