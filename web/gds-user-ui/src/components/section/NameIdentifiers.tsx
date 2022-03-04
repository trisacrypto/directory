import { HStack } from "@chakra-ui/react";
import Button from "components/ui/FormButton";
import FormLayout from "layouts/FormLayout";
import NameIdentifier from "./NameIdentifier";

const NameIdentifiers: React.FC = () => {
  return (
    <FormLayout>
      <NameIdentifier
        name="Name identifiers"
        description="The name and type of name by which the legal person is known."
      />
      <HStack width="100%" wrap="wrap" align="start" gap={4}>
        <Button>Add Legal Name</Button>
        <Button marginLeft="0 !important">Add Legal Name</Button>
        <Button marginLeft="0 !important">Add Legal Name</Button>
      </HStack>
    </FormLayout>
  );
};

export default NameIdentifiers;
