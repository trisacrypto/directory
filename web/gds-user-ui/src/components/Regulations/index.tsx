import { Grid, GridItem, VStack } from "@chakra-ui/react";
import DeleteButton from "components/ui/DeleteButton";
import FormButton from "components/ui/FormButton";
import InputFormControl from "components/ui/InputFormControl";

const Regulations: React.FC = () => {
  return (
    <VStack align="start" w="100%">
      {/* TODO: add a fieldArray to create new line */}
      <Grid templateColumns={{ base: "1fr auto" }} gap={6} width="100%">
        <GridItem>
          <InputFormControl controlId="applicable_regulation" />
        </GridItem>
        <GridItem display="flex" alignItems="center">
          <DeleteButton tooltip={{ label: "Remove line" }} />
        </GridItem>
      </Grid>
      <FormButton borderRadius={5}>Add Regulation</FormButton>
    </VStack>
  );
};

export default Regulations;
