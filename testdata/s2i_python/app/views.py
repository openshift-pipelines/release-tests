from django.shortcuts import render
from django.http import HttpResponse
from django.utils.timezone import now

def home(request):
    name = "RedHat"
    current_date = now().strftime("%Y-%m-%d %H:%M:%S")
    response = f"Name: {name}, Current Date: {current_date}"
    return HttpResponse(response)

# Create your views here.
