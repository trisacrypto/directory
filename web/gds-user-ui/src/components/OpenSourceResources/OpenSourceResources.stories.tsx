import { Meta, Story } from "@storybook/react";
import OpenSourceResources from ".";

export default {
  title: "components/OpenSourceRessources",
  component: OpenSourceResources,
} as Meta;

const Template: Story = (args) => <OpenSourceResources {...args} />;

export const Standard = Template.bind({});
Standard.args = {};
