import { Meta, Story } from "@storybook/react";
import DetailsCard from ".";

export default {
  title: "components/DetailsCard",
  component: DetailsCard,
} as Meta;

const Template: Story = (args) => <DetailsCard {...args} />;

export const Default = Template.bind({});
Default.args = {};
