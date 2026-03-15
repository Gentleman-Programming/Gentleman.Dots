---
name: dotnet
description: >
  .NET 9 / ASP.NET Core patterns with Minimal APIs, Clean Architecture, and EF Core.
  Trigger: When writing C# code, .NET APIs, or Entity Framework models.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## Minimal APIs (REQUIRED for new endpoints)

```csharp
// ✅ Minimal API with typed results
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddScoped<IOrderService, OrderService>();

var app = builder.Build();

var orders = app.MapGroup("/api/orders")
    .WithTags("Orders")
    .RequireAuthorization();

orders.MapGet("/", async (IOrderService service) =>
    TypedResults.Ok(await service.GetAllAsync()));

orders.MapGet("/{id:int}", async (int id, IOrderService service) =>
    await service.GetByIdAsync(id) is { } order
        ? TypedResults.Ok(order)
        : TypedResults.NotFound());

orders.MapPost("/", async (CreateOrderDto dto, IOrderService service) =>
{
    var order = await service.CreateAsync(dto);
    return TypedResults.Created($"/api/orders/{order.Id}", order);
});

app.Run();
```

## Primary Constructors (REQUIRED for DI)

```csharp
// ✅ Primary constructor for dependency injection
public class OrderService(
    AppDbContext db,
    ILogger<OrderService> logger,
    IMapper mapper) : IOrderService
{
    public async Task<List<OrderDto>> GetAllAsync()
    {
        logger.LogInformation("Fetching all orders");
        var orders = await db.Orders.AsNoTracking().ToListAsync();
        return mapper.Map<List<OrderDto>>(orders);
    }
}

// ❌ NEVER: Manual field assignment from constructor
public class OrderService
{
    private readonly AppDbContext _db;
    public OrderService(AppDbContext db) { _db = db; }
}
```

## Clean Architecture Layers

```
src/
├── Domain/              # Entities, value objects, interfaces
│   ├── Entities/
│   ├── Interfaces/      # IOrderRepository, IUnitOfWork
│   └── ValueObjects/
├── Application/         # Use cases, DTOs, validators
│   ├── Orders/
│   │   ├── Commands/    # CreateOrderCommand, handler
│   │   ├── Queries/     # GetOrderQuery, handler
│   │   └── Dtos/
│   └── Common/          # Behaviors, mappings
├── Infrastructure/      # EF Core, external services
│   ├── Persistence/     # DbContext, configs, migrations
│   ├── Services/        # Email, storage, etc.
│   └── DependencyInjection.cs
└── WebApi/              # Minimal API endpoints, middleware
    ├── Endpoints/
    └── Program.cs
```

## Entity Framework Core Patterns

```csharp
// ✅ Entity configuration (Fluent API, NOT data annotations)
public class OrderConfiguration : IEntityTypeConfiguration<Order>
{
    public void Configure(EntityTypeBuilder<Order> builder)
    {
        builder.HasKey(o => o.Id);
        builder.Property(o => o.Total).HasPrecision(18, 2);
        builder.HasMany(o => o.Items)
            .WithOne(i => i.Order)
            .HasForeignKey(i => i.OrderId)
            .OnDelete(DeleteBehavior.Cascade);
        builder.HasIndex(o => o.CreatedAt);
    }
}

// ✅ DbContext with configuration scanning
public class AppDbContext(DbContextOptions<AppDbContext> options)
    : DbContext(options)
{
    public DbSet<Order> Orders => Set<Order>();
    public DbSet<Product> Products => Set<Product>();

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.ApplyConfigurationsFromAssembly(
            typeof(AppDbContext).Assembly);
    }
}
```

## Repository Pattern (Optional)

```csharp
// ✅ Generic repository only if you need abstraction over EF
public interface IRepository<T> where T : class
{
    Task<T?> GetByIdAsync(int id);
    Task<List<T>> GetAllAsync();
    void Add(T entity);
    void Remove(T entity);
}

// ✅ Prefer using DbContext directly for simple cases
// Repositories add value only when you need to swap persistence
```

## DTOs and Mapping

```csharp
// ✅ Records for DTOs (immutable, concise)
public record CreateOrderDto(string CustomerName, List<OrderItemDto> Items);
public record OrderDto(int Id, string CustomerName, decimal Total, DateTime CreatedAt);
public record OrderItemDto(int ProductId, int Quantity);

// ❌ NEVER: Expose domain entities in API responses
```

## Result Pattern (No Exceptions for Control Flow)

```csharp
// ✅ Result type for expected failures
public class Result<T>
{
    public T? Value { get; }
    public string? Error { get; }
    public bool IsSuccess => Error is null;

    public static Result<T> Success(T value) => new() { Value = value };
    public static Result<T> Failure(string error) => new() { Error = error };
}

// ❌ NEVER: throw for business logic (e.g., "order not found")
```

## Dependency Injection Registration

```csharp
// ✅ Extension methods per layer
public static class DependencyInjection
{
    public static IServiceCollection AddInfrastructure(
        this IServiceCollection services, IConfiguration config)
    {
        services.AddDbContext<AppDbContext>(options =>
            options.UseSqlServer(config.GetConnectionString("Default")));

        services.AddScoped<IOrderRepository, OrderRepository>();
        services.AddScoped<IUnitOfWork, UnitOfWork>();

        return services;
    }
}

// In Program.cs:
builder.Services.AddInfrastructure(builder.Configuration);
```

## Global Error Handling

```csharp
// ✅ Problem Details with IExceptionHandler (.NET 8+)
public class GlobalExceptionHandler(ILogger<GlobalExceptionHandler> logger)
    : IExceptionHandler
{
    public async ValueTask<bool> TryHandleAsync(
        HttpContext context, Exception exception, CancellationToken ct)
    {
        logger.LogError(exception, "Unhandled exception");
        var problem = new ProblemDetails
        {
            Status = StatusCodes.Status500InternalServerError,
            Title = "Server Error",
        };
        context.Response.StatusCode = problem.Status.Value;
        await context.Response.WriteAsJsonAsync(problem, ct);
        return true;
    }
}
```

## Keywords
dotnet, .net, asp.net, c#, csharp, minimal api, ef core, entity framework, clean architecture, dependency injection
