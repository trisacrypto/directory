import { Heading, Text, Stack, VStack, Grid, GridItem, HStack, Box } from '@chakra-ui/react';
import Field from 'components/Field';
import DeleteButton from 'components/ui/DeleteButton';
import { getNameIdentiferTypeOptions } from 'constants/name-identifiers';
import React from 'react';
import {
  Control,
  FieldValues,
  RegisterOptions,
  useFieldArray,
  useFormContext
} from 'react-hook-form';

type NameIdentifierProps = {
  name: string;
  description: string;
  controlId?: string;
  register?: RegisterOptions;
  control?: Control<FieldValues, any>;
  heading?: string;
};

const NameIdentifier: React.ForwardRefExoticComponent<
  NameIdentifierProps & React.RefAttributes<unknown>
> = React.forwardRef((props, ref) => {
  const { name, controlId, description, heading } = props;
  const { register, control } = useFormContext();
  const { fields, remove, append } = useFieldArray({ name, control });

  React.useImperativeHandle(ref, () => ({
    addRow() {
      append({
        legal_person_name: '',
        legal_person_name_identifier_type: ''
      });
    }
  }));

  return (
    <Stack align="start" width="100%">
      {fields &&
        fields.map((field, index) => {
          return (
            <>
              {index < 1 && (
                <VStack align="start">
                  <Heading size="md">{heading}</Heading>
                  <Text size="sm">{description}</Text>
                </VStack>
              )}
              <HStack width="100%" key={field.id}>
                <Grid templateColumns={{ base: '1fr', md: 'repeat(2, 1fr)' }} gap={6} width="100%">
                  <GridItem>
                    <Field.Input
                      register={register}
                      name={`${name}[${index}].legal_person_name`}
                      controlId={controlId}
                    />
                  </GridItem>
                  <GridItem>
                    <Field.Select
                      register={register}
                      name={`${name}[${index}].legal_person_name_identifier_type`}
                      controlId={controlId}
                      options={getNameIdentiferTypeOptions()}
                    />
                  </GridItem>
                </Grid>
                <Box
                  paddingBottom={{ base: 2, md: 0 }}
                  alignSelf={{ base: 'flex-end', md: 'initial' }}>
                  <DeleteButton onDelete={() => remove(index)} tooltip={{ label: 'Remove line' }} />
                </Box>
              </HStack>
            </>
          );
        })}
    </Stack>
  );
});

NameIdentifier.displayName = 'NameIdentifier';

export default NameIdentifier;
