package vkcs

// func TestAccComputeV2InterfaceAttachImport_basic(t *testing.T) {
// 	resourceName := "openstack_compute_interface_attach_v2.ai_1"

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:          func() { testAccPreCheckCompute(t) },
// 		ProviderFactories: testAccProviders,
// 		CheckDestroy:      testAccCheckComputeInterfaceAttachDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccComputeInterfaceAttachBasic(),
// 			},
// 			{
// 				ResourceName:      resourceName,
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 				ImportStateVerifyIgnore: []string{
// 					"admin_pass",
// 				},
// 			},
// 		},
// 	})
// }
