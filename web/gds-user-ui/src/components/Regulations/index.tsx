import { Button, Grid, GridItem, VStack, Text } from '@chakra-ui/react';
import { t } from '@lingui/macro';
import { Trans } from '@lingui/react';
import DeleteButton from 'components/ui/DeleteButton';
import InputFormControl from 'components/ui/InputFormControl';
import { useFieldArray, useFormContext } from 'react-hook-form';

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
            <Text>{`${name}[${index}]`}</Text>
            <InputFormControl
              controlId="applicable_regulation"
              {...register(`${name}[${index}]`)}
            />
          </GridItem>
          <GridItem display="flex" alignItems="center">
            <DeleteButton onDelete={() => remove(index)} tooltip={{ label: t`Remove line` }} />
          </GridItem>
        </Grid>
      ))}
      <Button onClick={() => append({})} borderRadius={5}>
        <Trans id="Add Regulation">Add Regulation</Trans>
      </Button>
    </VStack>
  );
};

export default Regulations;
