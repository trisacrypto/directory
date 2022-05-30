import { Heading, Text, Stack, VStack, Grid, GridItem, HStack, Box } from '@chakra-ui/react';
import { t } from '@lingui/macro';
import DeleteButton from 'components/ui/DeleteButton';
import InputFormControl from 'components/ui/InputFormControl';
import SelectFormControl from 'components/ui/SelectFormControl';
import { getNameIdentiferTypeOptions } from 'constants/name-identifiers';
import React from 'react';
import {
  Control,
  Controller,
  FieldValues,
  RegisterOptions,
  useFieldArray,
  useFormContext
} from 'react-hook-form';
import { getValueByPathname } from 'utils/utils';

type NameIdentifierProps = {
  name: string;
  description: string;
  controlId?: string;
  register?: RegisterOptions;
  control?: Control<FieldValues, any>;
  heading?: string;
  type?: string;
};

const NameIdentifier: React.ForwardRefExoticComponent<
  NameIdentifierProps & React.RefAttributes<unknown>
> = React.forwardRef((props, ref) => {
  const {
    register,
    control,
    formState: { errors },
    getValues
  } = useFormContext();
  const { name, controlId, description, heading, type } = props;
  const { fields, remove, append } = useFieldArray({ name, control });

  const nameIdentiferTypeOptions = getNameIdentiferTypeOptions();
  React.useImperativeHandle(ref, () => ({
    addRow() {
      append({
        legal_person_name: '',
        legal_person_name_identifier_type: ''
      });
    }
  }));

  const getOrganizationNameValue = getValues('organization_name');

  const getOrganizationName = (index: number) => {
    if (type === 'legal' && index === 0) {
      return getOrganizationNameValue;
    }
  };

  return (
    <Stack align="start" width="100%">
      {fields &&
        fields.map((field, index) => {
          return (
            <React.Fragment key={field.id}>
              {index < 1 && (
                <VStack align="start">
                  <Heading size="md">{heading}</Heading>
                  <Text size="sm">{description}</Text>
                </VStack>
              )}
              <HStack width="100%" key={field.id}>
                <Grid templateColumns={{ base: '1fr', md: 'repeat(2, 1fr)' }} gap={6} width="100%">
                  <GridItem>
                    <InputFormControl
                      controlId={`${name}[${index}].legal_person_name`}
                      data-testid="legal_person_name"
                      placeholder={getOrganizationName(index) || ''}
                      isInvalid={getValueByPathname(errors, `${name}[${index}].legal_person_name`)}
                      formHelperText={
                        getValueByPathname(errors, `${name}[${index}].legal_person_name`)?.message
                      }
                      {...register(`${name}[${index}].legal_person_name`)}
                    />
                  </GridItem>
                  <GridItem>
                    <Controller
                      name={`${name}[${index}].legal_person_name_identifier_type`}
                      control={control}
                      render={({ field: f }) => (
                        <SelectFormControl
                          onBlur={f.onBlur}
                          controlId={controlId!}
                          data-testid="legal_person_name_identifier_type"
                          name={f.name}
                          ref={f.ref}
                          isDisabled={(index === 0 && type && type === 'legal') || false}
                          isInvalid={getValueByPathname(errors, f.name)}
                          formHelperText={getValueByPathname(errors, f.name)?.message}
                          formatOptionLabel={(data: any) => (
                            <>
                              {data.label} {t`Name`}
                            </>
                          )}
                          options={getNameIdentiferTypeOptions()}
                          onChange={(newValue: any) => f.onChange(newValue.value)}
                          value={nameIdentiferTypeOptions.find(
                            (option) => option.value === f.value
                          )}
                        />
                      )}
                    />
                  </GridItem>
                </Grid>

                <Box
                  paddingBottom={{ base: 2, md: 0 }}
                  alignSelf={{ base: 'flex-end', md: 'initial' }}>
                  <DeleteButton
                    onDelete={() => remove(index)}
                    tooltip={{ label: t`Remove line` }}
                    isDisabled={type === 'legal' && index === 0}
                  />
                </Box>
              </HStack>
            </React.Fragment>
          );
        })}
    </Stack>
  );
});

NameIdentifier.displayName = 'NameIdentifier';

export default NameIdentifier;
