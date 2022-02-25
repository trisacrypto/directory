import React from "react";
import {

  Flex,
  Heading,
  Stack,
  Text,
} from "@chakra-ui/react";

interface LandingHeaderProps {
  
  title: string,
  description?: string
}
// we should add props to the LandingHead component to allow it to update content dynamically
const LandingHead : React.FC<any> = ({title, description} : LandingHeaderProps): any => {
  return (
      <Flex
      bgGradient="linear(270deg,#24a9df,#1aebb4)"
      color="white"
      width="100%"
      justifyContent="center"
      alignItems="center"
      direction="column"
      padding={4}
    >
      <Stack
        textAlign={"center"}
        color="white"
        align={"center"}
        spacing={{ base: 8, md: 10 }}
        py={{ base: 10, md: 20 }}
        >
       <Heading
          fontWeight={600}
          fontSize={{ base: "3xl", sm: "4xl", md: "6xl" }}
          lineHeight={"110%"}
        >
       {title || 'TRISA Global Directory Service'}
        </Heading>
        {description ? <Text maxW={"2xl"}>{description}</Text> :
         (<Text maxW={"2xl"}>
            Become Travel Rule compliant. <br/>
            Apply to Become a TRISA certified Virtual Asset Service Provider.'
  
        </Text>)
        }
       
      </Stack>
    </Flex>
  );
};

export default LandingHead;
