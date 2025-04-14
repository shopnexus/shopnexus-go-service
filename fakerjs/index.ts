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
} from "../node_modules/.prisma/client";

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
  const createdAccounts = await prisma.account.findMany({
    where: {
      username: {
        in: accountsData.map((a) => a.username),
      },
    },
  });

  // Create related admin and user accounts individually since we need the base account IDs
  for (const account of createdAccounts) {
    if (account.role === Role.ADMIN) {
      await prisma.adminAccount.create({
        data: {
          id: account.id,
        },
      });
    } else {
      await prisma.userAccount.create({
        data: {
          id: account.id,
          email: faker.internet.email(),
          phone: faker.phone.number({ style: "international" }),
          gender: randomEnum<Gender>({
            MALE: Gender.MALE,
            FEMALE: Gender.FEMALE,
            OTHER: Gender.OTHER,
          }),
          full_name: faker.person.fullName(),
        },
      });
    }
  }

  return createdAccounts;
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
  userAccounts: Account[],
  count: number
) {
  const addressesData = Array.from({ length: count }, () => {
    const userAccount =
      userAccounts[Math.floor(Math.random() * userAccounts.length)];
    return {
      user_id: userAccount.id,
      address: faker.location.streetAddress(),
      city: faker.location.city(),
      province: faker.location.state(),
      country: faker.location.country(),
    };
  });

  await prisma.address.createMany({
    data: addressesData,
    skipDuplicates: true,
  });

  return await prisma.address.findMany({ take: count });
}

// Generate Carts and Items
async function createCarts(
  prisma: TxPrisma,
  userAccounts: Account[],
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
  const salesData: any[] = [];

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

    // Create 1-3 products per payment
    const productCount = Math.floor(Math.random() * 3) + 1;
    const selectedProducts = faker.helpers.arrayElements(
      products.filter((p) => p.sold < p.quantity),
      Math.min(productCount, products.filter((p) => p.sold < p.quantity).length)
    );

    for (const product of selectedProducts) {
      const quantity = BigInt(Math.floor(Math.random() * 3) + 1);
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
        // Get available serials for the product
        const availableSerials = Array.from({ length: Number(quantity) }, () =>
          faker.string.alphanumeric(10).toUpperCase()
        );

        for (const serialId of availableSerials) {
          productSerialOnProductOnPaymentData.push({
            product_on_payment_id: productOnPaymentId,
            product_serial_id: serialId,
          });
        }
      }
    }
  }

  // Create payments
  if (paymentsData.length > 0) {
    await prisma.payment.createMany({
      data: paymentsData,
      skipDuplicates: true,
    });
  }

  // Create VNPay records
  if (vnpayData.length > 0) {
    await prisma.paymentVnpay.createMany({
      data: vnpayData,
      skipDuplicates: true,
    });
  }

  // Create product on payment records
  if (productOnPaymentData.length > 0) {
    await prisma.productOnPayment.createMany({
      data: productOnPaymentData,
      skipDuplicates: true,
    });
  }

  // Create product serial on payment records
  if (productSerialOnProductOnPaymentData.length > 0) {
    await prisma.productSerialOnProductOnPayment.createMany({
      data: productSerialOnProductOnPaymentData,
      skipDuplicates: true,
    });
  }

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

  for (let i = 0; i < refundCount; i++) {
    const payment = successfulPayments[i];
    const refundMethod = faker.helpers.arrayElement(refundMethods);
    const status = faker.helpers.arrayElement(statuses);
    const productOnPayments = await prisma.productOnPayment.findMany({
      where: { payment_id: payment.id },
    });

    // Create refund with proper address handling
    const refundId = BigInt(i + 1);
    refundsData.push({
      id: refundId,
      product_on_payment_id:
        productOnPayments[Math.floor(Math.random() * productOnPayments.length)]
          .id,
      method: refundMethod,
      status: "PENDING", // Always start with PENDING status
      reason: faker.lorem.sentence(),
      // Only include address for PICK_UP method
      address: refundMethod === "PICK_UP" ? faker.location.streetAddress() : "",
      date_created: faker.date.recent(),
      date_updated: faker.date.recent(),
    });

    // Add 0-3 resources per refund
    const resourceCount = Math.floor(Math.random() * 4);
    for (let j = 0; j < resourceCount; j++) {
      resourcesData.push({
        owner_id: refundId,
        url: faker.image.urlPicsumPhotos({ width: 800, height: 600 }),
      });
    }
  }

  // Create refunds
  if (refundsData.length > 0) {
    await prisma.refund.createMany({
      data: refundsData,
      skipDuplicates: true,
    });
  }

  // Create resources for refunds
  if (resourcesData.length > 0) {
    await prisma.resource.createMany({
      data: resourcesData,
      skipDuplicates: true,
    });
  }

  return await prisma.refund.findMany({ take: count });
}

// Generate Comments and Resources
async function createComments(
  prisma: TxPrisma,
  accounts: Account[],
  count: number
) {
  const commentsData: Prisma.CommentCreateManyInput[] = [];
  const resourcesData: Prisma.ResourceCreateManyInput[] = [];
  const products = await prisma.product.findMany();
  const productModels = await prisma.productModel.findMany();
  const brands = await prisma.brand.findMany();

  const addComment = (commentId: bigint, destId: bigint, type: CommentType) => {
    const account = accounts[Math.floor(Math.random() * accounts.length)];
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
        owner_id: commentId,
        url: faker.image.urlPicsumPhotos({ width: 800, height: 600 }),
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
  const resourcesData: any[] = [];

  // Add resources for brands
  for (const brand of brands) {
    // Add 1-3 resources per brand
    const resourceCount = Math.floor(Math.random() * 3) + 1;
    for (let i = 0; i < resourceCount; i++) {
      resourcesData.push({
        owner_id: brand.id,
        url: faker.image.url({ width: 800, height: 600 }),
      });
    }
  }

  // Add resources for product models
  for (const productModel of productModels) {
    // Add 2-5 resources per product model
    const resourceCount = Math.floor(Math.random() * 4) + 2;
    for (let i = 0; i < resourceCount; i++) {
      resourcesData.push({
        owner_id: productModel.id,
        url: faker.image.urlPicsumPhotos({ width: 800, height: 600 }),
      });
    }
  }

  // Add resources for products
  for (const product of products) {
    // Add 1-3 resources per product
    const resourceCount = Math.floor(Math.random() * 3) + 1;
    for (let i = 0; i < resourceCount; i++) {
      resourcesData.push({
        owner_id: product.id,
        url: faker.image.urlPicsumPhotos({ width: 800, height: 600 }),
      });
    }
  }

  // Add resources for refunds
  for (const refund of refunds) {
    // Add 0-2 resources per refund (some refunds might not have resources)
    const resourceCount = Math.floor(Math.random() * 3);
    for (let i = 0; i < resourceCount; i++) {
      resourcesData.push({
        owner_id: refund.id,
        url: faker.image.urlPicsumPhotos({ width: 800, height: 600 }),
      });
    }
  }

  await prisma.resource.createMany({
    data: resourcesData,
    skipDuplicates: true,
  });

  return await prisma.resource.findMany({ take: resourcesData.length });
}

// Main seeding function
async function main() {
  const prisma = new PrismaClient();
  try {
    console.log("Starting to seed database...");

    await prisma.$transaction(async (tx) => {
      // Create base data
      await createRoles(tx);
      const accounts = await createAccounts(tx, 10);
      const userAccounts = accounts.filter((a) => a.role === "USER");
      const brands = await createBrands(tx, 5);
      const tags = await createTags(tx, 8);
      const productTypes = await createProductTypes(tx, 5);
      const productModels = await createProductModels(
        tx,
        brands,
        tags,
        productTypes,
        20
      );
      const products = await createProducts(tx, productModels, 50);

      // Create related data
      const addresses = await createAddresses(tx, userAccounts, 15);
      const carts = await createCarts(tx, userAccounts, products);
      const sales = await createSales(tx, productModels, tags, brands, 10);
      // const payments = await createPayments(tx, userAccounts, products, 15);
      // const refunds = await createRefunds(tx, payments, 5);
      const comments = await createComments(tx, accounts, 20);

      // Create resources for brands, productModels, products, and refunds
      const resources = await createResources(
        tx,
        brands,
        productModels,
        products,
        [],
        0
      );

      console.log("Seeding completed successfully!");
    });
  } catch (error) {
    console.error("Error seeding database:", error);
  } finally {
    await prisma.$disconnect();
  }
}

main();
