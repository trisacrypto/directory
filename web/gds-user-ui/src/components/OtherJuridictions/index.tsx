import { Box, Button, Grid, GridItem, HStack } from '@chakra-ui/react';
import DeleteButton from 'components/ui/DeleteButton';
import InputFormControl from 'components/ui/InputFormControl';
import SelectFormControl from 'components/ui/SelectFormControl';
import { getCountriesOptions } from 'constants/countries';
import { Controller, useFieldArray, useFormContext } from 'react-hook-form';
import { t } from '@lingui/macro';
import { Trans } from '@lingui/react';
const OtherJuridictions: React.FC<{ name: string }> = ({ name }) => {
  const { control, register } = useFormContext();
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
              <Controller
                control={control}
                name={`${name}[${index}].country`}
                render={({ field: f }) => (
                  <SelectFormControl
                    ref={f.ref}
                    name={f.name}
                    value={getCountriesOptions().find((option) => option.value === f.value)}
                    onChange={(newValue: any) => f.onChange(newValue.value)}
                    options={getCountriesOptions()}
                    label={t`National Jurisdiction`}
                    controlId="country"
                  />
                )}
              />
            </GridItem>
            <GridItem>
              <InputFormControl
                label={t`Regulator Name`}
                controlId="regulator_name"
                {...register(`${name}[${index}].regulator_name`)}
              />
            </GridItem>
          </Grid>
          <Box marginTop="23px" alignSelf={{ base: 'flex-end', md: 'initial' }}>
            <DeleteButton onDelete={() => remove(index)} tooltip={{ label: t`Remove line` }} />
          </Box>
        </HStack>
      ))}

      <Button onClick={handleAddJuridictionClick} borderRadius={5}>
        <Trans id="Add Jurisdiction">Add Jurisdiction</Trans>
      </Button>
    </>
  );
};

export default OtherJuridictions;
