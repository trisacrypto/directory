import { Grid, GridItem, VStack } from '@chakra-ui/react';
import { t } from '@lingui/macro';
import { Trans } from '@lingui/react';
import DeleteButton from 'components/ui/DeleteButton';
import FormButton from 'components/ui/FormButton';
import InputFormControl from 'components/ui/InputFormControl';
import { Control, useFieldArray, UseFormRegister, useFormContext } from 'react-hook-form';

type RegulationsProps = {
  name: string;
};

const Regulations: React.FC<RegulationsProps> = ({ name }) => {
  const { control, register } = useFormContext();
  const { fields, append, remove } = useFieldArray({
    name,
    control
  });
  return (
    <VStack align="start" w="100%">
      {fields.map((field, index) => (
        <Grid key={field.id} templateColumns={{ base: '1fr auto' }} gap={6} width="100%">
          <GridItem>
            <InputFormControl
              controlId="applicable_regulation"
              {...register(`${name}[${index}].name`)}
            />
          </GridItem>
          <GridItem display="flex" alignItems="center">
            <DeleteButton onDelete={() => remove(index)} tooltip={{ label: t`Remove line` }} />
          </GridItem>
        </Grid>
      ))}
      <FormButton onClick={() => append({ name: '' })} borderRadius={5}>
        <Trans id="Add Regulation">Add Regulation</Trans>
      </FormButton>
    </VStack>
  );
};

export default Regulations;
