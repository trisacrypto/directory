import { WarningIcon } from '@chakra-ui/icons';
import { Heading, Stack, Text } from '@chakra-ui/react';
import InputFormControl from 'components/ui/InputFormControl';
import React from 'react';
import { useFormContext } from 'react-hook-form';
import { getDomain } from 'utils/utils';
import { t } from '@lingui/macro';

type TrisaImplementationFormProps = {
  headerText: string;
  name: string;
  type: 'TestNet' | 'MainNet';
};

const env = {
  TestNet: 'testnet',
  MainNet: 'trisa'
};

const TrisaImplementationForm: React.FC<TrisaImplementationFormProps> = ({
  headerText,
  name,
  type
}) => {
  const {
    register,
    formState: { errors },
    watch,
    getValues
  } = useFormContext();
  const commonName = watch(`${name}.common_name`);
  const trisaEndpoint = watch(`${name}.endpoint`);
  const [commonNameWarning, setCommonNameWarning] = React.useState<string | undefined>('');

  React.useEffect(() => {
    const trisaEndpointUri = trisaEndpoint?.split(':')[0];

    const warningMessage =
      trisaEndpointUri === commonName
        ? undefined
        : t`Common name should match the TRISA endpoint without the port.`;
    setCommonNameWarning(warningMessage);
  }, [commonName, trisaEndpoint]);

  const getCommonNameFormHelperText = () => {
    if ((errors[name] as any)?.common_name) {
      return (errors[name] as any)?.common_name.message;
    }
    if (commonNameWarning) {
      return (
        <Text color="yellow.500">
          <WarningIcon /> {commonNameWarning}
        </Text>
      );
    }

    return t`The common name for the mTLS certificate. This should match the TRISA endpoint without the port in most cases.`;
  };

  const domain = getValues('website') && getDomain(getValues('website'));

  return (
    <Stack>
      <Heading size="md">{headerText}</Heading>
      <InputFormControl
        label={t`TRISA Endpoint`}
        placeholder={`${env[type]}.${domain}:443`}
        formHelperText={
          (errors[name] as any)?.endpoint
            ? (errors[name] as any)?.endpoint?.message
            : t`The address and port of the TRISA endpoint for partner VASPs to connect on via gRPC.`
        }
        controlId="trisaEndpoint"
        isInvalid={!!(errors[name] as any)?.endpoint}
        {...register(`${name}.endpoint`)}
      />

      <InputFormControl
        label={t`Certificate Common Name`}
        placeholder={`${env[type]}.${domain}`}
        isInvalid={!!(errors[name] as any)?.common_name}
        formHelperText={getCommonNameFormHelperText()}
        controlId="certificateCommonName"
        {...register(`${name}.common_name`)}
      />
    </Stack>
  );
};

export default TrisaImplementationForm;
