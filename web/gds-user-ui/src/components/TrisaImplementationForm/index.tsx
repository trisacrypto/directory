import { Heading } from '@chakra-ui/react';
import InputFormControl from 'components/ui/InputFormControl';
import FormLayout from 'layouts/FormLayout';
import { useFormContext } from 'react-hook-form';

type TrisaImplementationFormProps = {
  headerText: string;
  name: string;
};

const TrisaImplementationForm: React.FC<TrisaImplementationFormProps> = ({ headerText, name }) => {
  const { register, getValues } = useFormContext();
  console.log('[]', getValues());
  return (
    <FormLayout>
      <Heading size="md">{headerText}</Heading>
      <InputFormControl
        label="TRISA Endpoint"
        placeholder="trisa.example.com:443"
        formHelperText="The address and port of the TRISA endpoint for partner VASPs to connect on via gRPC."
        controlId="trisaEndpoint"
        {...register(`${name}.endpoint`)}
      />

      <InputFormControl
        label="Certificate Common Name"
        placeholder="trisa.example.com"
        formHelperText="The common name for the mTLS certificate. This should match the TRISA endpoint without the port in most cases."
        controlId="certificateCommonName"
        {...register(`${name}.common_name`)}
      />
    </FormLayout>
  );
};

export default TrisaImplementationForm;
