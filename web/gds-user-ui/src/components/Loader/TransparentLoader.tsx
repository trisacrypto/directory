import { Stack, Spinner, Flex, Box, Text, VStack, position } from '@chakra-ui/react';

interface LoaderProps {
  title?: string;
}
// This loader is for containers that are loading
const TransparentLoader: React.FC<LoaderProps> = (props) => {
  return (
    <Flex
      position={'fixed'}
      width={'100%'}
      left={0}
      top={0}
      right={0}
      bottom={0}
      backgroundColor={'rgba(255,255,255,0.7)'}
      zIndex={1}>
      <VStack spacing={4} m={'auto'} py={10}>
        <Spinner color="blue.500" size="xl" />
        <Text>{props?.title}</Text>
      </VStack>
    </Flex>
  );
};
TransparentLoader.defaultProps = {
  title: 'Session expired , reconnecting...'
};

export default TransparentLoader;
