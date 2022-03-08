import { VStack } from "@chakra-ui/react";
import InputFormControl from "components/ui/InputFormControl";
import SelectFormControl from "components/ui/SelectFormControl";

const BasicDetailsForm: React.FC = () => {
  return (
    <>
      <VStack spacing={4} w="100%">
        <InputFormControl
          controlId="website"
          label="Website"
          name="website"
          inputProps={{ placeholder: "VASP Holdings LLC" }}
        />

        <InputFormControl
          controlId="website"
          label="Date of Incorporation / Establishment"
          name="website"
          formHelperText=""
          inputProps={{ placeholder: "21/01/2021", type: "date" }}
        />

        <SelectFormControl
          label="Business Category"
          placeholder="Select business category"
          controlId="business-category"
        />

        <SelectFormControl
          label="VASP Category"
          placeholder="Select VASP category"
          controlId="vasp-category"
          isMulti
          formHelperText="Please select as many categories needed to represent the types of virtual asset services your organization provides."
        />
      </VStack>
    </>
  );
};

export default BasicDetailsForm;
