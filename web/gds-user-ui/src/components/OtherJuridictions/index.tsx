import { Box, Grid, GridItem, HStack } from '@chakra-ui/react';
import DeleteButton from 'components/ui/DeleteButton';
import FormButton from 'components/ui/FormButton';
import InputFormControl from 'components/ui/InputFormControl';
import SelectFormControl from 'components/ui/SelectFormControl';
import { getCountriesOptions } from 'constants/countries';
import { useFieldArray, useFormContext } from 'react-hook-form';

const OtherJuridictions: React.FC<{ name: string }> = ({ name }) => {
  const { control } = useFormContext();
  const { fields, append, remove } = useFieldArray({
    name,
    control
  });

  const handleAddJuridictionClick = () => {
    append({
      country: '',
      regulator_name: ''
    });
  };

  return (
    <>
      {fields.map((field, index) => (
        <HStack key={field.id}>
          <Grid templateColumns={{ base: '1fr 1fr', md: '2fr 1fr' }} gap={6} width="100%">
            <GridItem>
              <SelectFormControl
                options={getCountriesOptions()}
                label="National Jurisdiction"
                controlId="country"
              />
            </GridItem>
            <GridItem>
              <InputFormControl type="number" label="Regulator Name" controlId="regulator_name" />
            </GridItem>
          </Grid>
          <Box marginTop="23px" alignSelf={{ base: 'flex-end', md: 'initial' }}>
            <DeleteButton onDelete={() => remove(index)} tooltip={{ label: 'Remove line' }} />
          </Box>
        </HStack>
      ))}

      <FormButton onClick={handleAddJuridictionClick} borderRadius={5}>
        Add Jurisdiction
      </FormButton>
    </>
  );
};

export default OtherJuridictions;
