import { Grid, GridItem, VStack, Text } from "@chakra-ui/react";
import InputFormControl from "components/ui/InputFormControl";
import SelectFormControl from "components/ui/SelectFormControl";

const LegalPersonForm: React.FC<{}> = () => {
  return (
    <>
      <Text>Address 1</Text>
      <VStack spacing={3.5} align="start">
        <InputFormControl
          formHelperText="Address line 1 e.g. building name/number, street name"
          controlId="address_1"
        />

        <InputFormControl
          formHelperText="Address line 2 e.g. apartment or suite number"
          controlId="address_2"
        />

        <InputFormControl
          formHelperText="Address line 3 e.g. city, province, postal code"
          controlId="address_3"
        />

        <Grid
          templateColumns={{ base: "1fr", md: "repeat(2, 1fr)" }}
          gap={6}
          width="100%"
        >
          <GridItem>
            <InputFormControl formHelperText="Country" controlId="country" />
          </GridItem>
          <GridItem>
            <SelectFormControl
              formHelperText="Address Type"
              controlId="address_type"
            />
          </GridItem>
        </Grid>
      </VStack>
    </>
  );
};

export default LegalPersonForm;
