import { Heading, Link, Text } from "@chakra-ui/react";
import InputFormControl from "components/ui/InputFormControl";
import SelectFormControl from "components/ui/SelectFormControl";
import FormLayout from "layouts/FormLayout";

type NationalIdentificationProps = {};

const NationalIdentification: React.FC<NationalIdentificationProps> = () => {
  return (
    <FormLayout>
      <Heading size="md">National Identification</Heading>
      <Text>
        Please supply a valid national identification number. TRISA recommends
        the use of LEI numbers. For more information, please visit{" "}
        <Link href="https://gleif.org" color="blue.500" isExternal>
          GLEIF.org
        </Link>
      </Text>
      <InputFormControl
        label="Identification Number"
        controlId="identification_number"
        formHelperText="An identifier issued by an appropriate issuing authority"
      />
      <SelectFormControl
        label="Identification Type"
        controlId="identification_type"
      />
      <SelectFormControl
        label="Country of Issue"
        controlId="country_of_issue"
      />
      <InputFormControl
        label="Registration Authority"
        controlId="registration_authority"
        formHelperText="If the identifier is an LEI number, enter the ID used in the GLEIF Registration Authorities List."
      />
    </FormLayout>
  );
};

export default NationalIdentification;
