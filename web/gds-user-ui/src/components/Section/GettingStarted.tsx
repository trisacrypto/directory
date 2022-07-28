import * as React from 'react';
import { Stack, Text, VStack, Flex, Button } from '@chakra-ui/react';

const GettingStartedSection = () => {
  return (
    <Flex justifyItems={'center'} mx={{ lg: 70, md: 40, sm: 10 }} py={10}>
      <Stack direction={['column', 'row']} spacing={5} fontSize={{ lg: 24, sm: 18, md: 20 }}>
        <VStack textAlign={'center'} bg={'#E5EDF1'} p={5}>
          <Text> Step 1</Text>
          <Text pt={3} fontWeight={'bold'}>
            Create account
          </Text>
          <Text pb={3}>
            create your TRISA account with your VASP email address. Add collaborators in your
            organization.
          </Text>
          <Button bg={'white'} color={'#221F1F'}>
            create account
          </Button>
        </VStack>
        <VStack textAlign={'center'} bg={'#E5EDF1'} p={5}>
          <Text> Step 1</Text>
          <Text pt={3} fontWeight={'bold'}>
            Complete VASP Verification
          </Text>
          <Text pb={3}>
            Complete the multi-part TRISA verification form and due diligence process. Once
            approved, gain access to the Testnet.
          </Text>
          <Button bg={'white'} color={'#221F1F'}>
            Learn More
          </Button>
        </VStack>
        <VStack textAlign={'center'} bg={'#E5EDF1'} p={5}>
          <Text> Step 1</Text>
          <Text pt={3} fontWeight={'bold'}>
            Integrate and Comply
          </Text>
          <Text pb={3}>
            Set up your TRISA node or integrate with a 3rd-party Travel Rule solution. Complete
            testing and move to production.
          </Text>
          <Button bg={'white'} color={'#221F1F'}>
            Learn More
          </Button>
        </VStack>
      </Stack>
    </Flex>
  );
};

export default GettingStartedSection;
