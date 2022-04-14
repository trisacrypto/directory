import { Heading, Text, Stack, VStack, Grid, GridItem, HStack, Box } from '@chakra-ui/react';
import DeleteButton from 'components/ui/DeleteButton';
import InputFormControl from 'components/ui/InputFormControl';
import SelectFormControl from 'components/ui/SelectFormControl';
import { getNameIdentiferTypeOptions } from 'constants/name-identifiers';
import React, { useState, FC, useEffect } from 'react';
import {
  Control,
  Controller,
  FieldValues,
  RegisterOptions,
  useFieldArray,
  useFormContext
} from 'react-hook-form';
import { loadDefaultValueFromLocalStorage } from 'utils/localStorageHelper';
import { RootStateOrAny, useSelector } from 'react-redux';
import useCertificateStepper from 'hooks/useCertificateStepper';
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
    watch,
    setValue
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
  const getOrganizationName = watch('organization_name');
  const currentStep: number = useSelector((state: RootStateOrAny) => state.stepper.currentStep);
  useEffect(() => {
    if (currentStep === 2) {
      setValue(`entity.name.name_identifiers[0].legal_person_name`, getOrganizationName);
    }
  }, [getOrganizationName]);

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
                      // onValueChange={
                      //   index === 0 && getLegalNameDefaultValue(index, basicDetailOrganizationName)
                      // }
                      // isDisabled={(index === 0 && type && type === 'legal') || false}
                      {...register(`${name}[${index}].legal_person_name`)}
                    />
                  </GridItem>
                  <GridItem>
                    <Controller
                      name={`${name}[${index}].legal_person_name_identifier_type`}
                      control={control}
                      render={({ field: f }) => (
                        <SelectFormControl
                          controlId={controlId!}
                          name={f.name}
                          isDisabled={(index === 0 && type && type === 'legal') || false}
                          formatOptionLabel={(data: any) => <>{data.label} Name</>}
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
                    tooltip={{ label: 'Remove line' }}
                    isDisabled={index === 0}
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
