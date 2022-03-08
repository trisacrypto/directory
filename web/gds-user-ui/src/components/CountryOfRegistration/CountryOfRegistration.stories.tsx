import { Meta, Story } from "@storybook/react";
import CountryOfRegistration from ".";

type CountryOfRegistrationProps = {};

export default {
  title: "components/Country of Registration",
  component: CountryOfRegistration,
} as Meta<CountryOfRegistrationProps>;

const Template: Story<CountryOfRegistrationProps> = (args) => (
  <CountryOfRegistration {...args} />
);

export const Default = Template.bind({});
Default.args = {};
