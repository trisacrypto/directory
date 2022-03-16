import { Heading, Text, Stack, VStack, Grid, GridItem, HStack, Box } from '@chakra-ui/react';
import DeleteButton from 'components/ui/DeleteButton';
import InputFormControl from 'components/ui/InputFormControl';
import SelectFormControl from 'components/ui/SelectFormControl';
import React from 'react';

type NameIdentifierProps = {
  name: string;
  description: string;
};

const NameIdentifier: React.ForwardRefExoticComponent<
  NameIdentifierProps & React.RefAttributes<unknown>
> = React.forwardRef((props, ref) => {
  const { name, description } = props;
  React.useImperativeHandle(ref, () => ({
    addRow() {
      //    add function that append new row
    }
  }));

  return (
    <Stack align="start" width="100%">
      <VStack align="start">
        <Heading size="md">{name}</Heading>
        <Text size="sm">{description}</Text>
      </VStack>
      <HStack width="100%">
        <Grid templateColumns={{ base: '1fr', md: 'repeat(2, 1fr)' }} gap={6} width="100%">
          <GridItem>
            <InputFormControl name="legal_person_name" controlId="legal_person_name" />
          </GridItem>
          <GridItem>
            <SelectFormControl controlId="legal_person_name_identifier_type" />
          </GridItem>
        </Grid>
        <Box paddingBottom={{ base: 2, md: 0 }} alignSelf={{ base: 'flex-end', md: 'initial' }}>
          <DeleteButton tooltip={{ label: 'Remove line' }} />
        </Box>
      </HStack>
    </Stack>
  );
});

NameIdentifier.displayName = 'NameIdentifier';

export default NameIdentifier;
