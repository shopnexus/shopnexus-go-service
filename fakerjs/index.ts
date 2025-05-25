import { faker } from "@faker-js/faker";
import {
  CommentType,
  PrismaClient,
  Gender,
  Brand,
  Account,
  ProductModel,
  Product,
  PaymentMethod,
  Status,
  RefundMethod,
  Prisma,
  Tag,
  ProductType,
  Refund,
  ResourceType,
} from "../node_modules/.prisma/client";
import { UserAccount } from "@prisma/client";

type TxPrisma = Omit<
  PrismaClient,
  "$connect" | "$disconnect" | "$on" | "$transaction" | "$use" | "$extends"
>;

// Role
const Role = {
  ADMIN: "ADMIN",
  STAFF: "STAFF",
  USER: "USER",
};

// Helper function to generate random enum values
const randomEnum = <T>(enumObject: { [key: string]: T }): T => {
  const values = Object.values(enumObject);
  return values[Math.floor(Math.random() * values.length)];
};

async function createRoles(prisma: TxPrisma) {
  await prisma.role.createMany({
    data: [
      {
        name: Role.ADMIN,
      },
      {
        name: Role.STAFF,
      },
      {
        name: Role.USER,
      },
    ],
    skipDuplicates: true,
  });
}

async function createAccounts(prisma: TxPrisma, count: number) {
  // Generate Account with UserAccount or AdminAccount
  const accountsData: Prisma.AccountCreateManyInput[] = [];

  for (let i = 0; i < count; i++) {
    try {
      const isAdmin = Math.random() < 0.2; // 20% chance of being admin
      const username = faker.internet.username();

      accountsData.push({
        username,
        password: faker.internet.password(),
        role: isAdmin ? Role.ADMIN : Role.USER,
      });
    } catch (error) {
      if (
        !(
          error instanceof Prisma.PrismaClientKnownRequestError &&
          error.code === "P2002"
        )
      ) {
        throw error;
      }
      // Skip if unique constraint error (P2002)
      continue;
    }
  }

  await prisma.account.createMany({
    data: accountsData,
    skipDuplicates: true,
  });

  // Fetch created accounts to get their IDs
  const createdAccounts = await prisma.account.findMany({});

  // Prepare admin and user account data
  const adminAccountsData: Prisma.AdminAccountCreateManyInput[] = [];
  const userAccountsData: Prisma.UserAccountCreateManyInput[] = [];

  for (const account of createdAccounts) {
    if (account.role === Role.ADMIN) {
      adminAccountsData.push({
        id: account.id,
      });
    } else {
      userAccountsData.push({
        id: account.id,
        email: faker.internet.email(),
        phone: faker.phone.number({ style: "international" }),
        gender: randomEnum<Gender>({
          MALE: Gender.MALE,
          FEMALE: Gender.FEMALE,
          OTHER: Gender.OTHER,
        }),
        full_name: faker.person.fullName(),
        default_address_id: null,
      });
    }
  }

  // Create admin and user accounts in bulk
  if (adminAccountsData.length > 0) {
    await prisma.adminAccount.createMany({
      data: adminAccountsData,
      skipDuplicates: true,
    });
  }

  if (userAccountsData.length > 0) {
    await prisma.userAccount.createMany({
      data: userAccountsData,
      skipDuplicates: true,
    });
  }
}

// Generate Brands
async function createBrands(prisma: TxPrisma, count: number) {
  const brandsData = Array.from({ length: count }, () => ({
    name: faker.company.name(),
    description: faker.company.catchPhrase(),
  }));

  await prisma.brand.createMany({
    data: brandsData,
    skipDuplicates: true,
  });
  return await prisma.brand.findMany({ take: count });
}

// Generate Tags
async function createTags(prisma: TxPrisma, count: number) {
  const tagsData = Array.from({ length: count }, () => ({
    tag: faker.commerce.department(),
    description: faker.commerce.productDescription(),
  }));

  await prisma.tag.createMany({
    data: tagsData,
    skipDuplicates: true,
  });
  return await prisma.tag.findMany({ take: count });
}

// Generate ProductTypes
async function createProductTypes(prisma: TxPrisma, count: number) {
  const productTypesData = Array.from({ length: count }, () => ({
    name: faker.commerce.productMaterial(),
  }));

  await prisma.productType.createMany({
    data: productTypesData,
    skipDuplicates: true,
  });
  return await prisma.productType.findMany({ take: count });
}

// Generate ProductModels
async function createProductModels(
  prisma: TxPrisma,
  brands: Brand[],
  tags: Tag[],
  productTypes: ProductType[],
  count: number
) {
  const productModelsData: Prisma.ProductModelCreateManyInput[] = [];
  const tagOnProductData: Prisma.TagOnProductModelCreateManyInput[] = [];

  for (let i = 0; i < count; i++) {
    // Randomly assign 1-3 tags to each product model
    const tagCount = Math.floor(Math.random() * 3) + 1;
    const selectedTags = faker.helpers.arrayElements(tags, tagCount);
    const modelName = faker.commerce.productName();

    productModelsData.push({
      brand_id: brands[Math.floor(Math.random() * brands.length)].id,
      type: productTypes[Math.floor(Math.random() * productTypes.length)].id,
      name: modelName,
      description: faker.commerce.productDescription(),
      list_price: BigInt(
        parseInt(faker.commerce.price({ min: 100, max: 1000 })) * 1000
      ),
      date_manufactured: faker.date.past(),
    });

    // We'll need to fetch the created product models to get their IDs for the tag connections
  }

  await prisma.productModel.createMany({
    data: productModelsData,
    skipDuplicates: true,
  });

  // Fetch created product models to get their IDs
  const createdProductModels = await prisma.productModel.findMany({
    take: count,
  });

  // Create tag connections
  for (const productModel of createdProductModels) {
    const tagCount = Math.floor(Math.random() * 3) + 1;
    const selectedTags = faker.helpers.arrayElements(tags, tagCount);

    for (const tag of selectedTags) {
      tagOnProductData.push({
        product_model_id: productModel.id,
        tag: tag.tag,
      });
    }
  }

  if (tagOnProductData.length > 0) {
    await prisma.tagOnProductModel.createMany({
      data: tagOnProductData,
      skipDuplicates: true,
    });
  }

  return createdProductModels;
}

// Generate Products
async function createProducts(
  prisma: TxPrisma,
  productModels: ProductModel[],
  count: number
) {
  const productsData = Array.from({ length: count }, () => ({
    product_model_id:
      productModels[Math.floor(Math.random() * productModels.length)].id,
    quantity: BigInt(Math.floor(Math.random() * 100) + 10),
    sold: BigInt(Math.floor(Math.random() * 5)),
    add_price: BigInt(Math.floor(Math.random() * 50000)),
    is_active: Math.random() < 0.9, // 90% chance of being active
    can_combine: Math.random() < 0.3, // 30% chance of being combinable
    metadata: {
      color: faker.color.human(),
      size: Math.floor(Math.random() * 5) + 1,
    },
    date_created: faker.date.recent(),
    date_updated: faker.date.recent(),
  }));

  await prisma.product.createMany({
    data: productsData,
    skipDuplicates: true,
  });

  const createdProducts = await prisma.product.findMany({ take: count });

  // Create product serials for each product
  for (const product of createdProducts) {
    const serialsData = Array.from(
      { length: Number(product.quantity) },
      () => ({
        serial_id: faker.string.alphanumeric(10).toUpperCase(),
        product_id: product.id,
        is_sold: false,
        is_active: true,
        date_created: product.date_created,
        date_updated: product.date_updated,
      })
    );

    await prisma.productSerial.createMany({
      data: serialsData,
      skipDuplicates: true,
    });
  }

  return createdProducts;
}

// Generate Addresses
async function createAddresses(
  prisma: TxPrisma,
  userAccounts: UserAccount[],
  count: number
) {
  // Calculate how many addresses per user (rounded up)
  const addressesPerUser = Math.ceil(count / userAccounts.length);
  const addressesData: Prisma.AddressCreateManyInput[] = [];

  // Create addresses for each user
  for (const userAccount of userAccounts) {
    // Create 1 to addressesPerUser addresses for each user
    const userAddressCount = Math.floor(Math.random() * addressesPerUser) + 1;

    for (let i = 0; i < userAddressCount && addressesData.length < count; i++) {
      addressesData.push({
        user_id: userAccount.id,
        address: faker.location.streetAddress(),
        city: faker.location.city(),
        province: faker.location.state(),
        country: faker.location.country(),
        full_name: faker.person.fullName(),
        phone: faker.phone.number({ style: "international" }),
      });
    }
  }

  await prisma.address.createMany({
    data: addressesData,
    skipDuplicates: true,
  });

  return await prisma.address.findMany({ take: count });
}

// Generate Carts and Items
async function createCarts(
  prisma: TxPrisma,
  userAccounts: UserAccount[],
  products: Product[]
) {
  const cartsData: Prisma.CartCreateManyInput[] = [];
  const itemOnCartData: Prisma.ItemOnCartCreateManyInput[] = [];

  for (const userAccount of userAccounts) {
    cartsData.push({
      id: userAccount.id,
    });

    // Create cart items for each user
    const itemCount = Math.floor(Math.random() * 5) + 1;
    for (let i = 0; i < itemCount; i++) {
      itemOnCartData.push({
        cart_id: userAccount.id,
        product_id: products[Math.floor(Math.random() * products.length)].id,
        quantity: BigInt(Math.floor(Math.random() * 5) + 1),
      });
    }
  }

  await prisma.cart.createMany({
    data: cartsData,
    skipDuplicates: true,
  });

  if (itemOnCartData.length > 0) {
    await prisma.itemOnCart.createMany({
      data: itemOnCartData,
      skipDuplicates: true,
    });
  }

  return await prisma.cart.findMany({ take: userAccounts.length });
}

// Generate Sales
async function createSales(
  prisma: TxPrisma,
  productModels: ProductModel[],
  tags: any[],
  brands: Brand[],
  count: number
) {
  const salesData: Prisma.SaleCreateManyInput[] = [];

  for (let i = 0; i < count; i++) {
    const saleType = Math.random();
    const startDate = faker.date.recent();
    const endDate = faker.date.future({ refDate: startDate });
    const ran = Math.random();
    const discountPercent =
      ran < 0.7 ? Math.floor(Math.random() * 50) + 5 : null;
    const discountPrice =
      ran < 0.3 ? BigInt(Math.floor(Math.random() * 100000) + 10000) : null;

    salesData.push({
      tag:
        saleType < 0.33
          ? tags[Math.floor(Math.random() * tags.length)].tag
          : null,
      product_model_id:
        saleType >= 0.33 && saleType < 0.66
          ? productModels[Math.floor(Math.random() * productModels.length)].id
          : null,
      brand_id:
        saleType >= 0.66
          ? brands[Math.floor(Math.random() * brands.length)].id
          : null,
      date_started: startDate,
      date_ended: endDate,
      quantity: BigInt(Math.floor(Math.random() * 100) + 10),
      used: BigInt(Math.floor(Math.random() * 10)),
      is_active: true,
      discount_percent: discountPercent,
      discount_price: discountPrice,
      max_discount_price: BigInt(Math.floor(Math.random() * 500000) + 50000),
    });
  }

  await prisma.sale.createMany({
    data: salesData,
    skipDuplicates: true,
  });

  return await prisma.sale.findMany({ take: count });
}

// Generate Payments and ProductOnPayment
async function createPayments(
  prisma: TxPrisma,
  userAccounts: any[],
  products: Product[],
  count: number
) {
  const paymentsData: Prisma.PaymentCreateManyInput[] = [];
  const productOnPaymentData: Prisma.ProductOnPaymentCreateManyInput[] = [];
  const productSerialOnProductOnPaymentData: Prisma.ProductSerialOnProductOnPaymentCreateManyInput[] =
    [];
  const vnpayData: Prisma.PaymentVnpayCreateManyInput[] = [];
  const paymentMethods = Object.values(PaymentMethod);

  // First, fetch all available product serials
  const productSerials = await prisma.productSerial.findMany({
    where: {
      is_sold: false,
      is_active: true,
    },
  });

  // Group serials by product_id for easier access
  const serialsByProduct: Record<string, typeof productSerials> =
    productSerials.reduce((acc, serial) => {
      if (!acc[serial.product_id.toString()]) {
        acc[serial.product_id.toString()] = [];
      }
      acc[serial.product_id.toString()].push(serial);
      return acc;
    }, {} as Record<string, typeof productSerials>);

  for (let i = 0; i < count; i++) {
    const userAccount =
      userAccounts[Math.floor(Math.random() * userAccounts.length)];
    const method = faker.helpers.arrayElement(paymentMethods);
    const status: Status = faker.helpers.arrayElement(Object.values(Status));
    const totalPrice = BigInt(Math.floor(Math.random() * 1000000) + 100000);

    const paymentId = BigInt(i + 1);
    paymentsData.push({
      id: paymentId,
      user_id: userAccount.id,
      method: method,
      status: status,
      address: faker.location.streetAddress(),
      total: totalPrice,
      date_created: faker.date.recent(),
    });

    // Add VNPay data if payment method is VNPAY
    if (method === PaymentMethod.VNPAY) {
      vnpayData.push({
        id: paymentId,
        vnp_TxnRef: faker.string.alphanumeric(10),
        vnp_OrderInfo: `Payment for order ${paymentId}`,
        vnp_TransactionNo: faker.string.numeric(10),
        vnp_TransactionDate: faker.date.recent().toISOString(),
        vnp_CreateDate: faker.date.recent().toISOString(),
        vnp_IpAddr: faker.internet.ip(),
      });
    }

    // Filter products that have available serials
    const availableProducts = products.filter(
      (p) =>
        serialsByProduct[p.id.toString()] &&
        serialsByProduct[p.id.toString()].length > 0
    );

    if (availableProducts.length === 0) continue;

    // Create 1-3 products per payment
    const productCount = Math.floor(Math.random() * 3) + 1;
    const selectedProducts = faker.helpers.arrayElements(
      availableProducts,
      Math.min(productCount, availableProducts.length)
    );

    for (const product of selectedProducts) {
      const availableSerials = serialsByProduct[product.id.toString()];
      if (!availableSerials || availableSerials.length === 0) continue;

      const quantity = BigInt(
        Math.min(Math.floor(Math.random() * 3) + 1, availableSerials.length)
      );
      const price = BigInt(Math.floor(Math.random() * 500000) + 50000);
      const totalProductPrice = price * quantity;

      const productOnPaymentId = BigInt(productOnPaymentData.length + 1);
      productOnPaymentData.push({
        id: productOnPaymentId,
        payment_id: paymentId,
        product_id: product.id,
        quantity: quantity,
        price: price,
        total_price: totalProductPrice,
      });

      // Create serial numbers for successful payments
      if (status === Status.SUCCESS) {
        // Take the required number of serials and remove them from available serials
        const selectedSerials = availableSerials.splice(0, Number(quantity));

        for (const serial of selectedSerials) {
          productSerialOnProductOnPaymentData.push({
            product_on_payment_id: productOnPaymentId,
            product_serial_id: serial.serial_id,
          });

          // Mark serial as sold
          await prisma.productSerial.update({
            where: { serial_id: serial.serial_id },
            data: { is_sold: true },
          });
        }
      }
    }
  }

  // Bulk create all records
  if (paymentsData.length > 0) {
    await prisma.payment.createMany({
      data: paymentsData,
      skipDuplicates: true,
    });
  }

  if (vnpayData.length > 0) {
    await prisma.paymentVnpay.createMany({
      data: vnpayData,
      skipDuplicates: true,
    });
  }

  if (productOnPaymentData.length > 0) {
    await prisma.productOnPayment.createMany({
      data: productOnPaymentData,
      skipDuplicates: true,
    });
  }

  if (productSerialOnProductOnPaymentData.length > 0) {
    await prisma.productSerialOnProductOnPayment.createMany({
      data: productSerialOnProductOnPaymentData,
      skipDuplicates: true,
    });
  }

  // Set the next value for payment and product_on_payment sequences
  await prisma.$executeRawUnsafe(
    `ALTER SEQUENCE payment.base_id_seq RESTART WITH ${BigInt(count + 1)}`
  );
  await prisma.$executeRawUnsafe(
    `ALTER SEQUENCE payment.product_on_payment_id_seq RESTART WITH ${BigInt(
      productOnPaymentData.length + 1
    )}`
  );

  return await prisma.payment.findMany({ take: count });
}

// Generate Refunds
async function createRefunds(prisma: TxPrisma, payments: any[], count: number) {
  const refundsData: Prisma.RefundCreateManyInput[] = [];
  const resourcesData: Prisma.ResourceCreateManyInput[] = [];
  const refundMethods = Object.values(RefundMethod);
  const statuses = Object.values(Status);

  // Only create refunds for successful payments
  const successfulPayments = payments.filter((p) => p.status === "SUCCESS");
  const refundCount = Math.min(count, successfulPayments.length);

  // Get all product on payments first
  const productOnPayments = await prisma.productOnPayment.findMany({
    where: {
      payment_id: {
        in: successfulPayments.map((p) => p.id),
      },
    },
  });

  for (let i = 0; i < refundCount; i++) {
    const payment = successfulPayments[i];
    const refundMethod = faker.helpers.arrayElement(refundMethods);
    const status = faker.helpers.arrayElement(statuses);
    const availableProductOnPayments = productOnPayments.filter(
      (pop) => pop.payment_id === payment.id
    );

    if (availableProductOnPayments.length === 0) continue;

    const refundId = BigInt(i + 1);
    refundsData.push({
      id: refundId,
      product_on_payment_id:
        availableProductOnPayments[
          Math.floor(Math.random() * availableProductOnPayments.length)
        ].id,
      method: refundMethod,
      status: "PENDING",
      reason: faker.lorem.sentence(),
      address: refundMethod === "PICK_UP" ? faker.location.streetAddress() : "",
      date_created: faker.date.recent(),
      date_updated: faker.date.recent(),
    });

    // Add 0-3 resources per refund
    const resourceCount = Math.floor(Math.random() * 4);
    for (let j = 0; j < resourceCount; j++) {
      resourcesData.push({
        type: ResourceType.REFUND,
        owner_id: refundId,
        url: faker.image.urlPicsumPhotos({ width: 800, height: 600 }),
        order: j + 1,
      });
    }
  }

  // Create refunds and resources in bulk
  if (refundsData.length > 0) {
    await prisma.refund.createMany({
      data: refundsData,
      skipDuplicates: true,
    });
  }

  if (resourcesData.length > 0) {
    await prisma.resource.createMany({
      data: resourcesData,
      skipDuplicates: true,
    });
  }

  // Set the next value for the refund sequence
  await prisma.$executeRawUnsafe(
    `ALTER SEQUENCE payment.refund_id_seq RESTART WITH ${BigInt(
      refundCount + 1
    )}`
  );

  return await prisma.refund.findMany({ take: count });
}

// Generate Comments and Resources
async function createComments(
  prisma: TxPrisma,
  userAccounts: UserAccount[],
  count: number
) {
  const commentsData: Prisma.CommentCreateManyInput[] = [];
  const resourcesData: Prisma.ResourceCreateManyInput[] = [];
  const products = await prisma.product.findMany();
  const productModels = await prisma.productModel.findMany();
  const brands = await prisma.brand.findMany();

  const addComment = (commentId: bigint, destId: bigint, type: CommentType) => {
    const account =
      userAccounts[Math.floor(Math.random() * userAccounts.length)];
    commentsData.push({
      id: commentId,
      type: type,
      account_id: account.id,
      dest_id: destId,
      body: faker.lorem.paragraph(),
      upvote: BigInt(Math.floor(Math.random() * 100)),
      downvote: BigInt(Math.floor(Math.random() * 50)),
      score: Math.floor(Math.random() * 100) - 50,
      date_created: faker.date.past(),
      date_updated: faker.date.recent(),
    });

    // Add 0-3 resources per comment
    const resourceCount = Math.floor(Math.random() * 4);
    for (let j = 0; j < resourceCount; j++) {
      resourcesData.push({
        type: ResourceType.COMMENT,
        owner_id: commentId,
        url: faker.image.urlPicsumPhotos({ width: 800, height: 600 }),
        order: j + 1,
      });
    }
  };

  let commentId = BigInt(1);

  for (let i = 0; i < count; i++) {
    addComment(
      commentId++,
      productModels[Math.floor(Math.random() * productModels.length)].id,
      CommentType.PRODUCT_MODEL
    );
    addComment(
      commentId++,
      brands[Math.floor(Math.random() * brands.length)].id,
      CommentType.BRAND
    );
    addComment(
      commentId++,
      products[Math.floor(Math.random() * products.length)].id,
      CommentType.COMMENT
    );
  }

  // Set the next value for the comment sequence
  await prisma.$executeRawUnsafe(
    `ALTER SEQUENCE product.comment_id_seq RESTART WITH ${commentId}`
  );

  await prisma.comment.createMany({
    data: commentsData,
    skipDuplicates: true,
  });

  if (resourcesData.length > 0) {
    await prisma.resource.createMany({
      data: resourcesData,
      skipDuplicates: true,
    });
  }

  return await prisma.comment.findMany({ take: count });
}

// Generate Resources for brands, productModels, products, and refunds
async function createResources(
  prisma: TxPrisma,
  brands: Brand[],
  productModels: ProductModel[],
  products: Product[],
  refunds: Refund[],
  count: number
) {
  const resourcesData: Prisma.ResourceCreateManyInput[] = [];

  // Add resources for brands
  for (const brand of brands) {
    // Add 1-3 resources per brand
    const resourceCount = Math.floor(Math.random() * 3) + 1;
    for (let i = 0; i < resourceCount; i++) {
      resourcesData.push({
        type: ResourceType.BRAND,
        owner_id: brand.id,
        url: faker.image.url({ width: 800, height: 600 }),
        order: i + 1,
      });
    }
  }

  // Add resources for product models
  for (const productModel of productModels) {
    // Add 2-5 resources per product model
    const resourceCount = Math.floor(Math.random() * 4) + 2;
    for (let i = 0; i < resourceCount; i++) {
      resourcesData.push({
        type: ResourceType.PRODUCT_MODEL,
        owner_id: productModel.id,
        url: faker.image.urlPicsumPhotos({ width: 800, height: 600 }),
        order: i + 1,
      });
    }
  }

  // Add resources for products
  for (const product of products) {
    // Add 1-3 resources per product
    const resourceCount = Math.floor(Math.random() * 3) + 1;
    for (let i = 0; i < resourceCount; i++) {
      resourcesData.push({
        type: ResourceType.PRODUCT,
        owner_id: product.id,
        url: faker.image.urlPicsumPhotos({ width: 800, height: 600 }),
        order: i + 1,
      });
    }
  }

  // Add resources for refunds
  for (const refund of refunds) {
    // Add 0-2 resources per refund (some refunds might not have resources)
    const resourceCount = Math.floor(Math.random() * 3);
    for (let i = 0; i < resourceCount; i++) {
      resourcesData.push({
        type: ResourceType.REFUND,
        owner_id: refund.id,
        url: faker.image.urlPicsumPhotos({ width: 800, height: 600 }),
        order: i + 1,
      });
    }
  }

  await prisma.resource.createMany({
    data: resourcesData,
    skipDuplicates: true,
  });

  // Set the next value for the resource sequence
  await prisma.$executeRawUnsafe(
    `ALTER SEQUENCE product.resource_id_seq RESTART WITH ${BigInt(
      resourcesData.length + 1
    )}`
  );

  return await prisma.resource.findMany({ take: resourcesData.length });
}

// Define preset types and configurations
type PresetType = "light" | "medium" | "heavy";

interface SeedConfig {
  accounts: number;
  brands: number;
  tags: number;
  productTypes: number;
  productModels: number;
  products: number;
  addresses: number;
  sales: number;
  payments: number;
  comments: number;
}

// Update main function to accept preset parameter
async function main(preset: PresetType = "light") {
  const prisma = new PrismaClient();
  const config = SEED_PRESETS[preset];

  try {
    console.log(`Starting to seed database with ${preset} preset...`);

    await prisma.$transaction(
      async (tx) => {
        // Create base data
        await createRoles(tx);
        await createAccounts(tx, config.accounts);
        const userAccounts = await tx.userAccount.findMany();
        const brands = await createBrands(tx, config.brands);
        const tags = await createTags(tx, config.tags);
        const productTypes = await createProductTypes(tx, config.productTypes);
        const productModels = await createProductModels(
          tx,
          brands,
          tags,
          productTypes,
          config.productModels
        );
        const products = await createProducts(
          tx,
          productModels,
          config.products
        );

        // Create related data
        const addresses = await createAddresses(
          tx,
          userAccounts,
          config.addresses
        );
        const carts = await createCarts(tx, userAccounts, products);
        const sales = await createSales(
          tx,
          productModels,
          tags,
          brands,
          config.sales
        );
        const payments = await createPayments(
          tx,
          userAccounts,
          products,
          config.payments
        );
        const refunds = await createRefunds(
          tx,
          payments,
          Math.floor(config.payments * 0.1)
        ); // 10% of payments
        const comments = await createComments(
          tx,
          userAccounts,
          config.comments
        );

        // Create resources
        const resources = await createResources(
          tx,
          brands,
          productModels,
          products,
          refunds,
          Math.floor(config.productModels * 2) // Average 2 resources per product model
        );

        console.log(`Seeding completed successfully with ${preset} preset!`);
      },
      {
        timeout: 1000000000,
      }
    );
  } catch (error) {
    console.error("Error seeding database:", error);
  } finally {
    await prisma.$disconnect();
  }
}

const SEED_PRESETS: Record<PresetType, SeedConfig> = {
  light: {
    accounts: 10,
    brands: 5,
    tags: 8,
    productTypes: 5,
    productModels: 20,
    products: 50,
    addresses: 15,
    sales: 10,
    payments: 15,
    comments: 20,
  },
  medium: {
    accounts: 1000, // 1k users
    brands: 50, // 50 different brands
    tags: 100, // 100 different tags
    productTypes: 30, // 30 product types
    productModels: 500, // 500 different product models
    products: 2000, // 2k product variants
    addresses: 1500, // ~1.5 addresses per user
    sales: 200, // 200 different sales/promotions
    payments: 3000, // 3k orders
    comments: 5000, // 5k comments
  },
  heavy: {
    accounts: 10000, // 10k users
    brands: 200, // 200 different brands
    tags: 300, // 300 different tags
    productTypes: 100, // 100 product types
    productModels: 2000, // 2k different product models
    products: 10000, // 10k product variants
    addresses: 15000, // ~1.5 addresses per user
    sales: 1000, // 1k different sales/promotions
    payments: 30000, // 30k orders
    comments: 50000, // 50k comments
  },
};

// You can now call main with different presets
main("medium"); // or main('medium') or main('heavy')
