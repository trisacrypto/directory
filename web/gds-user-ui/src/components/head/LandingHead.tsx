import React from 'react';
import { Flex, Heading, Stack, Text, HStack, Button } from '@chakra-ui/react';

interface LandingHeaderProps {
  title: string;
  description?: string;
  hasBtn?: boolean;
}
// we should add props to the LandingHead component to allow it to update content dynamically
const LandingHead: React.FC<any> = ({ title, description, hasBtn }: LandingHeaderProps): any => {
  return (
    <Flex
      bgGradient="linear(270deg,#24a9df,#1aebb4)"
      color="white"
      width="100%"
      height={286}
      justifyContent="center"
      direction="column"
      padding={4}
      fontSize={'2xl'}>
      <Stack textAlign={'center'} color="white" spacing={{ base: 3 }} px={60}>
        <Heading fontWeight={600} fontSize={{ md: '4xl', sm: '2xl' }} lineHeight={'80%'}>
          {title || 'TRISA Global Directory Service'}
        </Heading>
        {description ? (
          <Text fontSize={{ base: 'xl', md: '2xl', sm: 'lg' }}>{description}</Text>
        ) : (
          <Text fontSize={{ base: '30px', md: '2xl', sm: 'lg' }}>
            Become Travel Rule compliant. <br />
            Apply to Become a TRISA certified Virtual Asset Service Provider.
          </Text>
        )}
        {hasBtn && (
          <>
            <Stack justifyContent={'center'} pt={4} spacing={4} direction={['column', 'row']}>
              <Button bg={'white'} color={'black'} p={2} minWidth={306} minHeight={65}>
                {' '}
                Why Join TRISA{' '}
              </Button>
              <Button bg={'white'} color={'black'} p={2} minWidth={306} minHeight={65}>
                {' '}
                Search Directory{' '}
              </Button>
            </Stack>
          </>
        )}
      </Stack>
    </Flex>
  );
};

export default LandingHead;
