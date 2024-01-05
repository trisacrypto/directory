import { Heading, Stack, Text } from '@chakra-ui/react';
import InputFormControl from 'components/ui/InputFormControl';
import PhoneNumberInput from 'components/ui/PhoneNumberInput';
import { Controller, useFormContext } from 'react-hook-form';
import get from 'lodash/get';
import { t } from '@lingui/macro';
import { Trans } from '@lingui/react';
type ContactFormProps = {
  title: string;
  description: string;
  name: string;
};

const ContactForm: React.FC<ContactFormProps> = ({ title, description, name }) => {
  const { register, control, formState } = useFormContext();
  const { errors } = formState;

  const getEmailInstruction = () => {
    if (name === 'contacts.legal') {
      return t`Please use the email address associated with your organization. Group or shared email addresses such as compliance@yourvasp.com are permitted if the email account is actively monitored.`;
    }
    return t`Please use the email address associated with your organization.`;
  };

  const getPhoneMessageHint = () => {
    if (name === 'contacts.legal') {
      return (
        <div data-testid="legal-contact-phone-number-hint">
          <Trans id="A business phone number is required to complete physical verification for MainNet registration. Please provide a phone number where the Legal/ Compliance contact can be contacted">
            A business phone number is required to complete physical verification for MainNet
            registration. Please provide a phone number where the Legal/ Compliance contact can be
            contacted
          </Trans>
          .
        </div>
      );
    }
    return (
      <div data-testid="legal-contact-phone-number-hint">
        <Trans id="If supplied, use full phone number with country code">
          If supplied, use full phone number with country code
        </Trans>
        .
      </div>
    );
  };

  return (
    <Stack>
      <Heading size="md" data-testid="title">
        {title}
      </Heading>
      <Text fontStyle="italic" data-testid="description">
        {description}
      </Text>
      <InputFormControl
        label={t`Full Name`}
        formHelperText={t`Preferred name for email communication.`}
        controlId="fullName"
        isInvalid={!!get(errors, `${name}.name`)}
        {...register(`${name}.name`)}
        data-testid="fullName"
      />

      <InputFormControl
        label={t`Email Address`}
        formHelperText={
          get(errors, `${name}.email`)
            ? get(errors, `${name}.email`)?.message
            : (getEmailInstruction() as string)
        }
        controlId="fullName"
        type="email"
        isInvalid={!!get(errors, `${name}.email`)}
        {...register(`${name}.email`)}
        data-testid="email"
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
              isInvalid={!!get(errors, `${name}.phone`)}
              label={t`Phone Number `}
              formHelperText={
                (get(errors, `${name}.phone`)
                  ? get(errors, `${name}.phone`)?.message
                  : getPhoneMessageHint()) as string
              }
              controlId="phoneNumber"
              data-testid="phoneNumber"
            />
          );
        }}
      />
    </Stack>
  );
};

export default ContactForm;
