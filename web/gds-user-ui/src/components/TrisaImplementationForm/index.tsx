import { WarningIcon } from '@chakra-ui/icons';
import { Heading, Text } from '@chakra-ui/react';
import InputFormControl from 'components/ui/InputFormControl';
import FormLayout from 'layouts/FormLayout';
import React from 'react';
import { useFormContext } from 'react-hook-form';

type TrisaImplementationFormProps = {
  headerText: string;
  name: string;
};

const TrisaImplementationForm: React.FC<TrisaImplementationFormProps> = ({ headerText, name }) => {
  const {
    register,
    formState: { errors },
    watch
  } = useFormContext();
  const commonName = watch(`${name}.common_name`);
  const trisaEndpoint = watch(`${name}.trisa_endpoint`);
  const [commonNameWarning, setCommonNameWarning] = React.useState<string | undefined>('');

  React.useEffect(() => {
    const trisaEndpointUri = trisaEndpoint.split(':')[0];
    const warningMessage =
      trisaEndpointUri === commonName
        ? undefined
        : 'common name should match the TRISA endpoint without the port';
    setCommonNameWarning(warningMessage);
  }, [commonName, trisaEndpoint]);

  return (
    <FormLayout>
      <Heading size="md">{headerText}</Heading>
      <InputFormControl
        label="TRISA Endpoint"
        placeholder="trisa.example.com:443"
        formHelperText={
          errors[name]?.endpoint
            ? errors[name]?.endpoint?.message
            : 'The address and port of the TRISA endpoint for partner VASPs to connect on via gRPC.'
        }
        controlId="trisaEndpoint"
        isInvalid={!!errors[name]?.endpoint}
        {...register(`${name}.endpoint`)}
      />

      <InputFormControl
        label="Certificate Common Name"
        placeholder="trisa.example.com"
        isInvalid={!!errors[name]?.common_name}
        formHelperText={
          commonNameWarning ? (
            <Text color="yellow.500">
              <WarningIcon /> {commonNameWarning}
            </Text>
          ) : (
            'The common name for the mTLS certificate. This should match the TRISA endpoint without the port in most cases.'
          )
        }
        controlId="certificateCommonName"
        {...register(`${name}.common_name`)}
      />
    </FormLayout>
  );
};

export default TrisaImplementationForm;
