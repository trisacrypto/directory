import { Heading } from '@chakra-ui/react';
import InputFormControl from 'components/ui/InputFormControl';
import FormLayout from 'layouts/FormLayout';

type TrisaImplementationFormProps = {
  headerText: string;
};

const TrisaImplementationForm: React.FC<TrisaImplementationFormProps> = ({ headerText }) => {
  return (
    <FormLayout>
      <Heading size="md">{headerText}</Heading>
      <InputFormControl
        label="TRISA Endpoint"
        placeholder="trisa.example.com:443"
        formHelperText="The address and port of the TRISA endpoint for partner VASPs to connect on via gRPC."
        controlId="trisaEndpoint"
      />

      <InputFormControl
        label="Certificate Common Name"
        placeholder="trisa.example.com"
        formHelperText="The common name for the mTLS certificate. This should match the TRISA endpoint without the port in most cases."
        controlId="certificateCommonName"
      />
    </FormLayout>
  );
};

export default TrisaImplementationForm;
